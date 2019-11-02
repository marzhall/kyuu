// Command jsonfs allows for the consumption and manipulation of a JSON
// object as a file system hierarchy.
package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"aqwari.net/net/styx"
)

var (
	addr         = flag.String("a", ":5640", "Port to listen on")
	debug        = flag.Bool("D", false, "trace 9P messages")
	verbose      = flag.Bool("v", false, "print extra info")
	max_messages = flag.Int("m", 20000, "maximum number of messages that can be held by a each queue")
)

type server struct {
	file map[string]interface{}
}

var logrequests styx.HandlerFunc = func(s *styx.Session) {
	for s.Next() {
		if *verbose {
			log.Printf("%q %T %s", s.User, s.Request(), s.Request().Path())
		}
	}
}

func main() {
	flag.Parse()
	log.SetPrefix("")
	log.SetFlags(0)

	var srv server
	srv.file = make(map[string]interface{})
	var styxServer styx.Server
	if *verbose {
		styxServer.ErrorLog = log.New(os.Stderr, "", 0)
	}
	if *debug {
		styxServer.TraceLog = log.New(os.Stderr, "", 0)
	}
	styxServer.Addr = *addr
	styxServer.Handler = styx.Stack(logrequests, &srv)

	log.Fatal(styxServer.ListenAndServe())
}

func walkTo(v interface{}, loc string) (interface{}, interface{}, bool) {
	cwd := v
	parts := strings.FieldsFunc(loc, func(r rune) bool { return r == '/' })
	var parent interface{}

	for _, p := range parts {
		switch v := cwd.(type) {
		case map[string]interface{}:
			parent = v
			if child, ok := v[p]; !ok {
				return nil, nil, false
			} else {
				cwd = child
			}
		case []interface{}:
			parent = v
			i, err := strconv.Atoi(p)
			if err != nil {
				return nil, nil, false
			}
			if len(v) <= i {
				return nil, nil, false
			}
			cwd = v[i]
		default:
			return nil, nil, false
		}
	}
	return parent, cwd, true
}

func (srv *server) Serve9P(s *styx.Session) {
	for s.Next() {
		t := s.Request()
		parent, file, ok := walkTo(srv.file, t.Path())
		if !ok {
			t.Rerror("no such file or directory")
			continue
		}
		fi := &stat{name: path.Base(t.Path()), file: &fakefile{v: file, name: t.Path()}}
		switch t := t.(type) {
		case styx.Twalk:
			t.Rwalk(fi, nil)
		case styx.Topen:
			switch v := file.(type) {
			case map[string]interface{}, []interface{}:
				t.Ropen(mkdir(v), nil)
			default:
				if _, ok := internal_queues[t.Path()]; !ok {
					internal_queues[t.Path()] = make(chan messageOrEOF, *max_messages)
				}
				t.Ropen(&fakefile{v: file, name: t.Path()}, nil)
			}
		case styx.Tutimes:
			t.Rutimes(nil)
		case styx.Ttruncate:
			t.Rtruncate(nil)
		/*case styx.Wstat:
		t.Rstat(fi, nil)*/
		case styx.Tstat:
			t.Rstat(fi, nil)
		case styx.Tcreate:
			switch v := file.(type) {
			case map[string]interface{}:
				if t.Mode.IsDir() {
					dir := make(map[string]interface{})
					v[t.Name] = dir
					t.Rcreate(mkdir(dir), nil)
				} else {
					v[t.Name] = new(bytes.Buffer)
					internal_queues[t.Path()] = make(chan messageOrEOF, *max_messages)
					t.Rcreate(&fakefile{
						v:    v[t.Name],
						name: t.Path(),
						set:  func(s string) { v[t.Name] = s },
					}, nil)
				}
			default:
				t.Rerror("%s is not a directory", t.Path())
			}
		case styx.Tremove:
			switch v := file.(type) {
			case map[string]interface{}:
				if len(v) > 0 {
					t.Rerror("directory is not empty")
					break
				}
				if parent != nil {
					if m, ok := parent.(map[string]interface{}); ok {
						delete(m, path.Base(t.Path()))
						t.Rremove(nil)
					} else {
						t.Rerror("cannot delete array element yet")
						break
					}
				} else {
					t.Rerror("permission denied")
				}
			default:
				if parent != nil {
					if m, ok := parent.(map[string]interface{}); ok {
						delete(m, path.Base(t.Path()))
						_, ok := internal_queues[t.Path()]
						if ok {
							delete(internal_queues, t.Path())
						}

						t.Rremove(nil)
					} else {
						t.Rerror("Can't rm {0}", t.Path())
					}
				} else {
					t.Rerror("permission denied")
				}
			}
		}
	}
}
