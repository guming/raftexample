package main

import (
	"sync"
	"encoding/gob"
	"bytes"
	"log"
)

type kvstore struct{
	proposeC    chan<- string
	mu          sync.RWMutex
	kvStore     map[string]string

}

type kv struct {
	Key string
	Val string
}

func newKVStore(proposeC chan<- string, commitC <-chan *string, errorC <-chan error) *kvstore {
	s:=&kvstore{proposeC:proposeC,kvStore:make(map[string]string)}
	go s.readCommits(commitC, errorC)
	return s
}

func (s *kvstore) readCommits(commitC <-chan *string, errorC <-chan error) {
	for data := range commitC {
		if data == nil {
			continue
		}

		var dataKv kv
		dec := gob.NewDecoder(bytes.NewBufferString(*data))
		if err := dec.Decode(&dataKv); err != nil {
			log.Fatalf("rafte xample: could not decode message (%v)", err)
		}
		s.mu.Lock()
		s.kvStore[dataKv.Key] = dataKv.Val
		s.mu.Unlock()
	}
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}

func (s *kvstore) Lookup(key string) (string, bool) {
	s.mu.RLock()
	v, ok := s.kvStore[key]
	s.mu.RUnlock()
	return v, ok
}

func (s *kvstore) Propose(k string, v string) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(kv{k, v}); err != nil {
		log.Fatal(err)
	}
	s.proposeC <- buf.String()
}