package main

/* notes:
want to create a tree hierarchy mirroring the hierarchy of the
queues; then, want to create one or more files for each queue
in which to store records of both additions and deliveries of
messages ot it.

One restart, files should be able to be read in an arbitrary order.
Messages are kept track of by UUID.
*/

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

type MessagePersister struct {
	queued   *os.File
	dequeued *os.File
}

var Persister *MessagePersister = nil

func MakeMessagePersister(storagePath string) *MessagePersister {
	if storagePath == "" {
		err := error(nil)
		storagePath, err = os.Getwd()
		if err != nil {
			log.Fatalf("Didn't get a storage path, and can't use the current directory. Err is \n%s", err)
		}

		storagePath = path.Join(storagePath, "kyuu_storage")
		err = os.MkdirAll(storagePath, 0755)
		if err != nil {
			log.Fatalf("Couldn't create our storage directory. Err is \n%s", err)
		}
	}

	queuedPath := path.Join(storagePath, "queued")
	dequeuedPath := path.Join(storagePath, "dequeued")
	queued, err := os.OpenFile(queuedPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	dequeued, err := os.OpenFile(dequeuedPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}

	return &MessagePersister{queued, dequeued}
}

func (m *MessagePersister) Shutdown() {
	m.queued.Close()
	m.dequeued.Close()
}

func (m *MessagePersister) MessageQueued(msg []byte) error {
	marshall, err := json.Marshal(msg)
	if err != nil {
		log.Println("error marshalling ", msg)
		return err
	}

	_, err = m.queued.Write(marshall)
	if err != nil {
		log.Println("Error writing ", msg)
		return err
	}

	_, err = m.queued.Write([]byte("\n"))
	if err != nil {
		log.Println("Error writing endline for ", msg)
		return err
	}

	return nil
}

func (m *MessagePersister) MessageDelivered(msg []byte) error {
	marshall, err := json.Marshal(msg)
	if err != nil {
		log.Println("error marshalling ", msg)
		return err
	}

	_, err = m.dequeued.Write(marshall)
	if err != nil {
		log.Println("Error writing ", msg)
		return err
	}

	_, err = m.dequeued.Write([]byte("\n"))
	if err != nil {
		log.Println("Error writing endline for ", msg)
		return err
	}

	return nil
}
