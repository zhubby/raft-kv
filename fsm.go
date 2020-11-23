package raft_kv

import (
	"bytes"
	"encoding/gob"
	"github.com/hashicorp/raft"
	"io"
	"io/ioutil"
)

type fsm struct {
	kv *memoryKV
}

type logData struct {
	Name  string
	Value string
}

func (f *fsm) Apply(raftLog *raft.Log) interface{} {
	data := &logData{}
	if err := gob.NewDecoder(bytes.NewBuffer(raftLog.Data)).Decode(data); err != nil {
		return err
	}
	f.kv.Set(data.Name, data.Value)
	return nil
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	return snapshot{f.kv}, nil
}

func (f *fsm) Restore(read io.ReadCloser) error {
	b, err := ioutil.ReadAll(read)
	if err != nil {
		return err
	}
	return f.kv.Unmarshal(b)
}
