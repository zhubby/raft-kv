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
	defer s.RUnlock()
	return s.data[key]
}

func (s *memoryKV) Set(key, value string) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = value
}

func (s *memoryKV) Del(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, key)
}

func (s *memoryKV) Marshal() ([]byte, error) {
	s.RLocker()
	defer s.RUnlock()
	var b = bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(s.data)
	return b.Bytes(), err
}

func (s *memoryKV) Unmarshal(b []byte) error {
	s.Lock()
	defer s.Unlock()
	return gob.NewDecoder(bytes.NewBuffer(b)).Decode(s)
}
