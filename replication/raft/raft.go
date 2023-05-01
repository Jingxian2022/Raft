package raft

import (
	"context"
	"encoding/json"
	"log"
	"modist/orchestrator/node"
	pb "modist/proto"
	"sync"
	"time"
)

// RETRIES is the number of retries upon failure. By default, we have 3
const RETRIES = 3

// RETRY_TIME is the amount of time to wait between retries
const RETRY_TIME = 100 * time.Millisecond

type State struct {
	// In-memory key value store
	store map[string]string
	mu    sync.RWMutex

	// Channels given back by the underlying Raft implementation
	proposeC chan<- []byte
	commitC  <-chan *commit
	// TODO: what are the differences???????????

	proposedKVs      map[KV]struct{}
	replicateSuccess chan bool

	// Observability
	log *log.Logger

	// The public-facing API that this replicator must implement
	pb.ReplicatorServer
}

type Args struct {
	Node   *node.Node
	Config *Config
}

type KV struct {
	Key   string
	Value string
}

// Configure is called by the orchestrator to start this node
//
// The "args" are any to support any replicator that might need arbitrary
// set of configuration values.
//
// You are free to add any fields to the State struct that you feel may help
// with implementing the ReplicateKey or GetReplicatedKey functions. Likewise,
// you may call any additional functions following the second TODO. Apart
// from the two TODOs, DO NOT EDIT any other part of this function.
func Configure(args any) *State {
	a := args.(Args)

	node := a.Node
	config := a.Config

	proposeC := make(chan []byte)
	commitC := NewRaftNode(node, config, proposeC)

	s := &State{
		store:    make(map[string]string),
		proposeC: proposeC,
		commitC:  commitC,

		log: node.Log,

		// TODO(students): [Raft] Initialize any additional fields and add to State struct
		proposedKVs: make(map[KV]struct{}),
		//FIXME:????? do i need mutex?
		replicateSuccess: make(chan bool),
	}

	// We registered RaftRPCServer when calling NewRaftNode above so we only need to
	// register ReplicatorServer here
	s.log.Printf("starting gRPC server at %s", node.Addr.Host)
	grpcServer := node.GrpcServer
	pb.RegisterReplicatorServer(grpcServer, s)
	go grpcServer.Serve(node.Listener)

	// TODO(students): [Raft] Call helper functions if needed
	go s.commitCListener()

	return s
}

func (s *State) commitCListener() {
	for {
		select {
		case commit := <-s.commitC:
			if commit == nil {
				return
			}
			cmt := CommitMsg{} // TODO: what should be the type of commit?
			err := json.Unmarshal(*commit, &cmt)
			if err != nil {
				s.log.Printf("error unmarshalling commit: %v", err)
				continue
			}
			kv := KV{Key: cmt.key, Value: cmt.value}
			if _, ok := s.proposedKVs[kv]; ok && cmt.success {
				s.replicateSuccess <- true
				delete(s.proposedKVs, kv)
			}
			s.mu.Lock()
			s.store[cmt.key] = cmt.value
			s.mu.Unlock()
		}
	}
}

// ReplicateKey replicates the (key, value) given in the PutRequest by relaying it to the
// underlying Raft cluster/implementation.
//
// You should first marshal the KV struct using the encoding/json library. Then, put the
// marshalled struct in the proposeC channel and wait until the Raft cluster has confirmed that
// it has been committed. Think about how you can use the commitC channel to find out when
// something has been committed.
//
// If the proposal has not been committed after RETRY_TIME, you should retry it. You should retry
// the proposal a maximum of RETRIES times, at which point you can return a nil reply and an error.
func (s *State) ReplicateKey(ctx context.Context, r *pb.PutRequest) (*pb.PutReply, error) {

	// TODO(students): [Raft] Implement me!
	kvR, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	for i := 0; i < RETRIES; i++ {
		s.proposeC <- kvR
		s.proposedKVs[KV{Key: r.Key, Value: r.Value}] = struct{}{}
		select {
		case <-s.replicateSuccess:
			// unmarshal
			return &pb.PutReply{}, nil
		case <-time.After(RETRY_TIME):
			continue
		}
	}

	return nil, err
}

// GetReplicatedKey reads the given key from s.store. The implementation of
// this method should be pretty short (roughly 8-12 lines).
func (s *State) GetReplicatedKey(ctx context.Context, r *pb.GetRequest) (*pb.GetReply, error) {

	// TODO(students): [Raft] Implement me!
	s.mu.RLock()
	defer s.mu.RUnlock()
	if val, ok := s.store[r.Key]; ok {
		return &pb.GetReply{Value: val}, nil
	}
	return nil, nil
}
