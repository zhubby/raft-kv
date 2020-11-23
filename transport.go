package raft_kv

import (
	"github.com/hashicorp/raft"
	"net"
	"os"
	"time"
)

func newRaftTransport(opt *Options) (*raft.NetworkTransport, error) {
	tcp, err := net.ResolveTCPAddr("tcp", opt.RaftPeer)
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransport(tcp.String(), tcp, 5, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}
	return transport, nil
}
