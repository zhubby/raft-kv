package raft_kv

import (
	"bytes"
	"encoding/gob"
	"sync"
)

type KV interface {
	Get(key string) string
	Set(key, value string)
	Del(key string)
}

type memoryKV struct {
	data map[string]string
	sync.RWMutex
}

func newMemoryKV() *memoryKV {
	var m = make(map[string]string)
	return &memoryKV{data: m}
}

func (s *memoryKV) Get(key string) string {
	s.RLocker()
	value := s.data[key]
	s.RUnlock()
	return value
}

func (s *memoryKV) Set(key, value string) {
	s.Lock()
	s.data[key] = value
	s.Unlock()
}

func (s *memoryKV) Del(key string) {
	s.Lock()
	delete(s.data, key)
	s.Unlock()

}

func (s *memoryKV) Marshal() ([]byte, error) {
	var b = bytes.NewBuffer(nil)
	s.RLocker()
	err := gob.NewEncoder(b).Encode(s.data)
	s.RUnlock()
	return b.Bytes(), err
}

func (s *memoryKV) Unmarshal(b []byte) error {
	s.Lock()
	err := gob.NewDecoder(bytes.NewBuffer(b)).Decode(s)
	s.Unlock()
	return err
}
