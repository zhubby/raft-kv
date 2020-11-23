package raft_kv

import (
	"flag"
	"os"
)

var (
	dataDir   string
	raftPeer  string
	bootstrap bool
	joinPeer  string
)

type Options struct {
	DataDir   string
	RaftPeer  string
	Bootstrap bool
	JoinPeer  string
}

func Default() *Options {
	return &Options{
		DataDir:   dataDir,
		RaftPeer:  raftPeer,
		Bootstrap: bootstrap,
		JoinPeer:  joinPeer,
	}
}

func init() {
	flag.StringVar(&dataDir, "data-dir", os.TempDir(), "data dir")
	flag.StringVar(&raftPeer, "raft-peer", "127.0.0.1:8000", "raft master peer address")
	flag.StringVar(&joinPeer, "join-peer", "", "join peer address")
	flag.BoolVar(&bootstrap, "bootstrap", false, "start as master")
}
