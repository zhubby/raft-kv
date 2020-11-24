package raft_kv_test

import (
	"flag"
	raft_kv "github.com/zhubby/raft-kv"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Test_Node1(t *testing.T) {
	flag.Set("data-dir", filepath.Join(os.TempDir(), "node1"))
	flag.Set("bootstrap", "true")
	flag.Set("join-peer", "127.0.0.1:8001,127.0.0.1:8002")
	tx := time.NewTicker(5 * time.Second)

	kv, err := raft_kv.New(raft_kv.Default())
	if err != nil {
		t.Log(err)
	}

	for {
		select {
		case _ = <-tx.C:
			kv.Set(time.Now().String(), time.Now().String())
		}
	}
}

func Test_Node2(t *testing.T) {
	flag.Set("raft-peer", "127.0.0.1:8001")
	flag.Set("data-dir", filepath.Join(os.TempDir(), "node2"))
	tx := time.NewTicker(5 * time.Second)
	_, err := raft_kv.New(raft_kv.Default())
	if err != nil {
		t.Log(err)
	}
	for {
		select {
		case _ = <-tx.C:

		}
	}
}

func Test_Node3(t *testing.T) {
	flag.Set("raft-peer", "127.0.0.1:8002")
	flag.Set("data-dir", filepath.Join(os.TempDir(), "node3"))
	tx := time.NewTicker(5 * time.Second)
	_, err := raft_kv.New(raft_kv.Default())
	if err != nil {
		t.Log(err)
	}
	for {
		select {
		case _ = <-tx.C:

		}
	}
}

func init() {
	println(os.TempDir())
	os.MkdirAll(filepath.Join(os.TempDir(), "node1"), 0700)
	os.MkdirAll(filepath.Join(os.TempDir(), "node2"), 0700)
	os.MkdirAll(filepath.Join(os.TempDir(), "node3"), 0700)
}
