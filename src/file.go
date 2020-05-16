package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

var internal_queues = make(map[string]chan messageOrEOF)

// Turn Go types into files

type messageOrEOF struct {
	m     []byte
	isEof bool
}

type fakefile struct {
	v       interface{}
	offset  int64
	sendEof bool
	name    string
	set     func(s string)
}

func (f *fakefile) ReadAt(p []byte, off int64) (int, error) {
	if f.sendEof {
		f.sendEof = false
		return 0, io.EOF
	}

	select {
	case temp := <-internal_queues[f.name]:
		n := copy(p, temp.m)
		Persister.MessageDelivered(temp.m)
		f.sendEof = true
		return n, io.EOF
	default:
		return 0, io.EOF
	}
}

func (f *fakefile) Write(p []byte) (int, error) {
	tmp := make([]byte, len(p))
	copy(tmp, p)
	select {
	case internal_queues[f.name] <- messageOrEOF{m: tmp, isEof: false}:
		Persister.MessageQueued(tmp)
		return len(p), nil
	default:
		return 0, io.ErrShortWrite
	}
}

func (f *fakefile) WriteAt(p []byte, off int64) (int, error) {
	return f.Write(p)
}

func (f *fakefile) Close() error {
	if f.set != nil {
		f.set(fmt.Sprint(f.v))
	}
	return nil
}

func (f *fakefile) size() int64 {
	return 0
}

type stat struct {
	name string
	file *fakefile
}

func (s *stat) Name() string     { return s.name }
func (s *stat) Sys() interface{} { return s.file }

func (s *stat) ModTime() time.Time {
	return time.Now().Truncate(time.Hour)
}

func (s *stat) IsDir() bool {
	return s.Mode().IsDir()
}

func (s *stat) Mode() os.FileMode {
	switch s.file.v.(type) {
	case map[string]interface{}:
		return os.ModeDir | 0755
	case []interface{}:
		return os.ModeDir | 0755
	}
	return 0644
}

func (s *stat) Size() int64 {
	return s.file.size()
}

type dir struct {
	c    chan stat
	done chan struct{}
}

func mkdir(val interface{}) *dir {
	c := make(chan stat, 10)
	done := make(chan struct{})
	go func() {
		if m, ok := val.(map[string]interface{}); ok {
		LoopMap:
			for name, v := range m {
				select {
				case c <- stat{name: name, file: &fakefile{v: v}}:
				case <-done:
					break LoopMap
				}
			}
		} else if a, ok := val.([]interface{}); ok {
		LoopArray:
			for i, v := range a {
				name := strconv.Itoa(i)
				select {
				case c <- stat{name: name, file: &fakefile{v: v}}:
				case <-done:
					break LoopArray
				}
			}
		}
		close(c)
	}()
	return &dir{
		c:    c,
		done: done,
	}
}

func (d *dir) Readdir(n int) ([]os.FileInfo, error) {
	var err error
	fi := make([]os.FileInfo, 0, 10)
	for i := 0; i < n; i++ {
		s, ok := <-d.c
		if !ok {
			err = io.EOF
			break
		}
		fi = append(fi, &s)
	}
	return fi, err
}

func (d *dir) Close() error {
	close(d.done)
	return nil
}
