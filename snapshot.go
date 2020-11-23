package raft_kv

import (
	"fmt"
	"github.com/hashicorp/raft"
)

type snapshot struct {
	*memoryKV
}

func (s snapshot) Persist(sink raft.SnapshotSink) error {
	defer sink.Cancel()
	b, err := s.Marshal()
	if err != nil {
		return err
	}
	if _, err := sink.Write(b); err != nil {
		return err
	}
	if err := sink.Close(); err != nil {
		return err
	}
	return nil
}

func (s snapshot) Release() {
	fmt.Print("snapshot completed")
}
