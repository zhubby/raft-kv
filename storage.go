package raft_kv

import (
	"bytes"
	"encoding/gob"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"io"
	"sync"
)

type storage struct {
	data map[string]string
	sync.RWMutex
}

func (s *storage) Get(key string) string {
	s.RLocker()
	defer s.RUnlock()
	return s.data[key]
}

func (s *storage) Set(key, value string) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = value
}

func (s *storage) Marshal() ([]byte, error) {
	s.RLocker()
	defer s.RUnlock()
	var b = bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(s.data)
	return b.Bytes(), err
}

func (s *storage) Unmarshal(b []byte) error {
	s.Lock()
	defer s.Unlock()
	return gob.NewDecoder(bytes.NewBuffer(b)).Decode(s)
}

func New() {
	conf := raft.DefaultConfig()

	logStore, err := raftboltdb.NewBoltStore("")
	if err != nil {
		panic(err)
	}
	stableStore, err := raftboltdb.NewBoltStore("")
	if err != nil {
		panic(err)
	}

	raft.NewRaft(conf)
}

type logData struct {
	Name  string
	Value string
}

type fsm struct {
}

type fsmSnapshot struct {
	*storage
}

func (f fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	b, err := f.Marshal()
	sink.Write(f.Marshal())
	panic("implement me")
}

func (f fsmSnapshot) Release() {

}

func (f fsm) Apply(raftLog *raft.Log) interface{} {
	data := &logData{}
	if err := gob.NewDecoder(bytes.NewBuffer(raftLog.Data)).Decode(data); err != nil {
		panic(err)
	}
	return nil
}

func (f fsm) Snapshot() (raft.FSMSnapshot, error) {
	panic("implement me")
}

func (f fsm) Restore(io.ReadCloser) error {
	panic("implement me")
}
