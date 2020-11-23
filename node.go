package raft_kv

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RaftKV struct {
	r *raft.Raft
	*memoryKV
	opt    *Options
	logger hclog.Logger
}

func (kv RaftKV) join(peer string) error {
	future := kv.r.AddVoter(raft.ServerID(peer), raft.ServerAddress(peer), 0, 0)
	return future.Error()
}

func New(opt *Options) (*RaftKV, error) {
	var rkv = &RaftKV{}
	rkv.logger = hclog.New(&hclog.LoggerOptions{
		Name:   "raft-kv",
		Level:  hclog.LevelFromString("DEBUG"),
		Output: os.Stderr,
	})
	rkv.logger.Info("raft kv options", hclog.Fmt("%+v", opt))
	if opt == nil {
		rkv.opt = Default()
	} else {
		rkv.opt = opt
	}
	rkv.memoryKV = newMemoryKV()
	conf := raft.DefaultConfig()
	conf.LocalID = raft.ServerID(rkv.opt.RaftPeer)
	conf.NotifyCh = make(chan bool, 1)
	conf.SnapshotInterval = 30 * time.Second
	conf.SnapshotThreshold = 1
	dir := rkv.opt.DataDir
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(dir, "raft-log.bolt"))
	if err != nil {
		return nil, err
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(dir, "raft-stable.bolt"))
	if err != nil {
		return nil, err
	}
	snapshotStore, err := raft.NewFileSnapshotStore(dir, 1, os.Stderr)
	if err != nil {
		return nil, err
	}
	transport, err := newRaftTransport(rkv.opt)
	if err != nil {
		return nil, err
	}
	r, err := raft.NewRaft(conf, &fsm{}, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return nil, err
	}
	if rkv.opt.Bootstrap {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      conf.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		r.BootstrapCluster(configuration)
	}
	rkv.r = r
	go func() {
		for {
			select {
			case leader := <-rkv.r.LeaderCh():
				if leader {
					if rkv.opt.JoinPeer != "" {
						peers := strings.Split(rkv.opt.JoinPeer, ",")
						for _, v := range peers {
							if err := rkv.join(v); err != nil {
								rkv.logger.Error("join cluster failed", err)
							}
						}
					}
				}
			}
		}
	}()

	return rkv, nil
}
