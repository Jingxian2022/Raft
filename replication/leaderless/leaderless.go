package leaderless

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"modist/orchestrator/node"
	pb "modist/proto"
	"modist/replication/conflict"
	"modist/store"
	"sync"
	"time"
)

// Leaderless replication is a strategy where any node in a cluster can accept writes or reads,
// and this node is responsible for replicating to any other nodes it wants to. For fault
// tolerance, a minimum number of successful reads and writes can be required before returning,
// specified by the parameters R and W, respectively.

// safelyUpdateKey writes a new KV pair if it is not too old. It returns whether the key was
// updated and the most up-to-date key-value that this node knows about (which may be the new KV).
//
// There are 4 cases that you need to handle here. Understanding concurrent events in the context
// of the HappensBefore relation will be vital. Specifically, when updating a key, we need to
// make sure that the new key is as up-to-date as what we have stored locally (if applicable).
// This requires doing a get to check up-to-dateness, and then doing a put if what we have is,
// in fact, more up-to-date. If two operations are concurrent, use the conflict resolver in the
// state struct to resolve them and then update the state accordingly.
//
// Because these read and write operations need to be done atomically, we wrap them in a
// transaction.
// You will need the functions s.localStore.BeginTx() to create a tx,
// and then tx.Get() or tx.Put() when you want to read or write a key-value pair respectively.
// When you are done with your transaction, remember to commit the transaction with tx.Commit().
//
// See store/store.go for more transaction usage details.
func (s *State[T]) safelyUpdateKey(newKV *conflict.KV[T]) (updated bool, mostUpToDateKV *conflict.KV[T], err error) {

	// TODO(students): [Leaderless] Implement me!
	tx := s.localStore.BeginTx(false)
	defer tx.Commit()

	curKV, ok := tx.Get(newKV.Key)

	if !ok {
		tx.Put(newKV.Key, newKV)
		return true, newKV, nil
	}

	curClock := curKV.Clock
	newClock := newKV.Clock

	// no need to update
	if newClock.HappensBefore(curClock) {
		return false, curKV, nil
	}

	if curClock.HappensBefore(newClock) {
		tx.Put(newKV.Key, newKV)
		return true, newKV, nil
	}

	// conflicts
	updatedKV, err := s.conflictResolver.ResolveConcurrentEvents(curKV, newKV)

	if err != nil {
		return false, nil, err
	}

	tx.Put(newKV.Key, updatedKV)

	return updatedKV.Value == newKV.Value, updatedKV, nil
}

// getUpToDateKV returns the KV associated with the key from the local store, but only if the one
// that exists is at least as up-to-date as (or concurrent with) minimumClock. It returns the KV
// (if it passes this clock constraint) and a boolean denoting whether the key exists (regardless
// of the corresponding clock).
//
// Think about whether you need to use a transaction here.
//
// There are 4 cases to handle here, the same as in safelyUpdateKey.
func (s *State[T]) getUpToDateKV(key string, minimumClock T) (kv *conflict.KV[T], found bool) {

	// TODO(students): [Leaderless] Implement me!
	curKV, ok := s.localStore.Get(key)

	// cannot find
	if !ok {
		return nil, false
	} else if curKV.Clock.HappensBefore(minimumClock) {
		return nil, true
	}

	return curKV, true
}

// HandlePeerWrite attempts to write a KV being replicated from a peer node (not the client).
//
// Specifically, it performs a local write (if it passes up-to-date checks), returning whether
// the write was accepted and, if the write was not accepted, the most up-to-date key-value that
// this node knows about.
//
// safelyUpdateKey should handle most of this logic for you, but you'll need to assemble the
// HandlePeerWriteReply based off its return values -- i.e. the ResolvableKv represents the most
// up-to-date key-value.
//
// You'll notice that HandlePeerWriteReply's ResolvableKv is of type *pb.ResolvableKV. To convert
// the key-value returned from safelyUpdateKey to this type, simply use the .Proto() function of
// that key-value.
func (s *State[T]) HandlePeerWrite(ctx context.Context, r *pb.ResolvableKV) (*pb.HandlePeerWriteReply, error) {
	s.onMessageReceived(conflict.ClockFromProto[T](r.GetClock()))

	newKV := conflict.KVFromProto[T](r)
	s.log.Printf("HandlePeerWrite: received direct replication of %v", newKV)

	// TODO(students): [Leaderless] Implement me!
	updated, updatedKV, err := s.safelyUpdateKey(newKV)

	if err != nil {
		return nil, err
	}

	reply := pb.HandlePeerWriteReply{}
	reply.Accepted = updated
	reply.ResolvableKv = updatedKV.Proto()

	return &reply, nil
}

// replicateToNode performs a remote write of the given KV to the specified node, with 3 retries.
//
// Specifically, given the replica node id, send it an RPC to directly write the given key-value
// pair. You should create the RPC client using s.node.PeerConns and
// pb.NewBasicLeaderlessReplicatorClient.
//
// [IMPORTANT]: Don't forget to call s.onMessageSend() before sending the RPC: this will update
// the clock state as necessary. In addition, the function withRetries may also be helpful.
//
// Reference the handout for how to send an RPC to a given node using its ID! Also remember that
// you can convert a key-value into its protobuf counterpart (*pb.ResolvableKV) using the
// .Proto() function.
//
// There is one edge case you must consider: what should happen if the replica node does not
// accept your write?
func (s *State[T]) replicateToNode(ctx context.Context, kv *conflict.KV[T], replicaNodeID uint64) error {
	s.log.Printf("write to node being called for node %d", replicaNodeID)
	s.log.Printf("current node: %d", s.node.ID)

	// TODO(students): [Leaderless] Implement me!
	conn := s.node.PeerConns[uint64(replicaNodeID)]
	client := pb.NewBasicLeaderlessReplicatorClient(conn)

	err := s.withRetries(func() error {
		s.onMessageSend()
		_, err := client.HandlePeerWrite(ctx, kv.Proto())
		return err
	}, 3)

	return err
}

// ReplicateKey replicates the given key to W arbitrary nodes (one of which is the current node).
// The clock corresponding to this key should be returned back to the client in the reply.
//
// The write quorum W is defined for you in the state struct.
// You should also use the provided helper method dispatchToPeers to send direct replication
// RPCs to other replicas in parallel. You can define an anonymous function and pass it to
// dispatchToPeers (what function that you've written does it need to call?).
//
// The implementation of this method should be 7-10 lines.
func (s *State[T]) ReplicateKey(ctx context.Context, kv *pb.PutRequest) (*pb.PutReply, error) {
	// If the client didn't provide a clock, this must be their first request. Give them a new
	// clock starting now. Note that we do this same check in GetReplicatedKey.
	s.ensureClock(&kv.Clock)

	clientClock := conflict.ClockFromProto[T](kv.GetClock())
	s.onMessageReceived(clientClock)

	newKV := conflict.KVFromParts(kv.Key, kv.Value, s.conflictResolver.NewClock())
	s.log.Printf("ReplicateKey: called with KV %s", newKV)

	// TODO(students): [Leaderless] Implement me!
	err := s.dispatchToPeers(ctx, s.W, func(ctx context.Context, replicaNodeID uint64) error {
		return s.replicateToNode(ctx, newKV, replicaNodeID)
	})
	if err != nil {
		return nil, err
	}

	reply := pb.PutReply{}
	reply.Clock = kv.GetClock()

	return &reply, nil
}

// HandlePeerRead attempts to service a peer's read request by returning the KV from the current
// node's local store. Note that the read request will have a clock specifying the client's clock
// for this key. We return {Found: false} if the key-value we find locally is less up-to-date than
// the client's clock.
//
// getUpToDateKV should handle most of this logic for you, but you'll need to assemble the
// HandlePeerReadReply based off of its return values. As with HandlePeerWrite, you will want to
// convert the node's key-value type to the *pb.ResolvableKV type required by the
// HandlePeerReadReply.

// Because these read and write operations need to be done atomically, we wrap them in a
// transaction.
// You will need the functions s.localStore.BeginTx() to create a tx,
// and then tx.Get() or tx.Put() when you want to read or write a key-value pair respectively.
// When you are done with your transaction, remember to commit the transaction with tx.Commit().
//
// See store/store.go for more transaction usage details.
func (s *State[T]) HandlePeerRead(ctx context.Context, request *pb.Key) (*pb.HandlePeerReadReply, error) {
	requestKey := request.GetKey()
	requestClock := conflict.ClockFromProto[T](request.GetClock())
	s.onMessageReceived(requestClock) // clock: onMsgRcv

	s.log.Printf("Node %d's HandlePeerRead: received request for key %s", s.node.ID, requestKey)
	// TODO(students): [Leaderless] Implement me!
	tx := s.localStore.BeginTx(true)
	defer tx.Commit()

	reply := pb.HandlePeerReadReply{}

	kv, found := s.getUpToDateKV(requestKey, requestClock)

	reply.Found = found
	if found {
		reply.ResolvableKv = kv.Proto()
	}

	return &reply, nil
}

// readFromNode performs a remote read from the specified node, with 3 retries.
//
// Specifically, given the replica node id, send it an RPC to directly read the given key.
// Don't forget to call s.onMessageSend() before sending the RPC: this will update the clock
// state as necessary. Again, you may find withRetries useful.

// You should create the RPC client using s.node.PeerConns and pb.NewBasicLeaderlessReplicatorClient.

// An error should ONLY be returned if the actual RPC fails. If the replica node just does not
// have the specified key, do not return an error (you can think of this replica as just being
// really behind).
//
// Additionally, you will also want to convert the Proto representation of the key-value back to
// our node's key-value type (the inverse of what you did in HandlePeerRead/Write). To do this,
// you will want to use the conflict.KVFromProto function found in the conflict module.

// Reference the handout for how to send an RPC to a given node using its ID! Also remember that
// you can convert a key-value into its protobuf counterpart (*pb.ResolvableKV) using the
// .Proto() function.
func (s *State[T]) readFromNode(ctx context.Context, key string, replicaNodeID uint64, clientClock T) (*conflict.KV[T], error) {
	s.log.Printf("read from node being called for node %d", replicaNodeID)
	s.log.Printf("current node: %d", s.node.ID)
	// TODO(students): [Leaderless] Implement me!

	structedInKey := &pb.Key{Key: key, Clock: clientClock.Proto()}

	clientConn := s.node.PeerConns[replicaNodeID]
	c := pb.NewBasicLeaderlessReplicatorClient(clientConn)

	readReply := pb.HandlePeerReadReply{}

	retryFunc := func() error {
		s.onMessageSend()
		reply, err := c.HandlePeerRead(ctx, structedInKey)
		if err == nil {
			readReply.Found = reply.Found
			readReply.ResolvableKv = reply.GetResolvableKv()
		}
		return err
	}

	err := s.withRetries(retryFunc, 3)
	if err != nil {
		return nil, errors.New("HandlePeerRead error")
	}
	//if readReply == nil {
	//	return new(conflict.KV[T]), nil
	//}
	myKv := readReply.GetResolvableKv()
	if myKv == nil {
		return nil, nil
	}

	s.log.Printf("result of HandlePeerRead: %s", myKv.Value)
	return &conflict.KV[T]{Key: myKv.Key, Value: myKv.Value, Clock: clientClock}, nil

}

// PerformReadRepair performs synchronous read repair using the most up-to-date key-value pair,
// and a mapping of key-values that all the other nodes have. Up-to-dateness is checked, and if a
// node is not as up-to-date as the latest KV pair, it is read-repaired. Read repair of different
// nodes should be performed in parallel, and this function should block until all repairs are
// complete.
//
// Don't forget to call s.onMessageSend() before sending any RPCs. For simplicity, when repairing,
// you can call HandlePeerWrite without wrapping it with withRetries.
//
// The following requirements are a bit implementation-specific to this implementation of
// leaderless replication, but:
//  1. latestKV should be as up to date (if not more) than every KV in kvPairs
//  2. kvPairs can have values that are nil, meaning their associated node didn't have our key
func (s *State[T]) PerformReadRepair(ctx context.Context, latestKV *conflict.KV[T], kvPairs map[uint64]*conflict.KV[T]) {

	// TODO(students): [Leaderless] Implement me!
	var wg sync.WaitGroup
	for replicaNodeID, kv := range kvPairs {

		if !kv.Equals(latestKV) {
			s.log.Printf("updating node %d", replicaNodeID)
			wg.Add(1)
			go func() {
				defer wg.Done()
				clientConn := s.node.PeerConns[replicaNodeID]
				c := pb.NewBasicLeaderlessReplicatorClient(clientConn)
				s.onMessageSend()

				mykv := latestKV.Proto()
				// c.HandlePeerWrite(ctx, mykv)

				c.HandlePeerWrite(ctx, mykv)
			}()
		}
	}
	wg.Wait()
}

// GetReplicatedKey performs a quorum read of the system, also performing read repair.
//
// Like before, the read quorum R is defined for you in the state struct.
//
// You should use dispatchToPeers, passing a function that calls readFromNode. Once you've read
// from a node, keep track of the most up-to-date replied KV and store all replied KVs
// (this "keep track" code must be thread-safe!). To compute the most up-to-date KV, you need to
// consider and resolve KVs from other nodes that are concurrent with our current most up-to-date
// KV.
//
// Once you have the most up-to-date KV and the replies of the other nodes, you can
// PerformReadRepair. You should check the signature of PerformReadRepair to make sure you are
// computing its arguments in the right way.
func (s *State[T]) GetReplicatedKey(ctx context.Context, r *pb.GetRequest) (*pb.GetReply, error) {
	// If the client didn't provide a clock, this must be their first request. Give them a new
	// clock starting now. Note that we do this same check in ReplicateKey.
	s.ensureClock(&r.GetMetadata().Clock)

	clientClock := conflict.ClockFromProto[T](r.GetMetadata().GetClock())
	s.onMessageReceived(clientClock)

	s.log.Printf("GetReplicatedKey: key %s with clock %v", r.Key, clientClock)

	// TODO(students): [Leaderless] Implement me!
	resolver := s.conflictResolver
	R := s.R
	key := r.Key

	KVMap := make(map[uint64]*conflict.KV[T])
	var mutex = &sync.RWMutex{}
	// KVMap := M{mymap: make(map[uint64]*conflict.KV[T])}
	readFromNodeFunc := func(ctx context.Context, replicaNodeID uint64) error {
		getKV, err := s.readFromNode(ctx, key, replicaNodeID, clientClock)
		// getKV may equal to nil

		if err != nil {
			return err
		}

		mutex.Lock()
		KVMap[replicaNodeID] = getKV
		if getKV != nil {
			s.log.Printf("GetReplicatedKey: get value %s from node %d", getKV.Value, replicaNodeID)
		} else {
			s.log.Printf("GetReplicatedKey: get value nil from node %d", replicaNodeID)
		}
		mutex.Unlock()

		return nil
	}
	err := s.dispatchToPeers(ctx, R, readFromNodeFunc) //parallel
	if err != nil {
		return new(pb.GetReply), errors.New("GetReplicatedKey error")
	}

	// used to resolve conflict
	var kvs []*conflict.KV[T]
	for _, kv := range KVMap {
		// KVMap may had nil values
		if kv != nil {
			kvs = append(kvs, kv)
		}
	}

	if len(kvs) == 0 {
		return nil, errors.New("reading non-existent key")
	}

	latestKV, err := resolver.ResolveConcurrentEvents(kvs...)
	if err == nil {
		s.log.Printf("performing read repair...")
		s.PerformReadRepair(ctx, latestKV, KVMap)
	} else {
		return new(pb.GetReply), errors.New("No KV is read")
	}
	reply := pb.GetReply{Value: latestKV.Value, Clock: clientClock.Proto()}
	return &reply, nil
}

// ======================================
// DO NOT CHANGE ANY CODE BELOW THIS LINE
// ======================================

type State[T conflict.Clock] struct {
	// The node that this replicator server is part of
	node *node.Node

	// A centralized clock generator and conflict resolver
	conflictResolver conflict.Resolver[T]

	W int
	R int

	// The storage layer for local key-value pairs
	// Implemented in store/memory.go and store/store.go
	localStore store.Store[*conflict.KV[T]]

	// Lifecycle methods
	onReplicatorStart func(node *node.Node)
	onMessageSend     func()
	onMessageReceived func(clock T)

	// Determines what replicas to read/write a key(/value) from. This could be key-sensitive in
	// the future, but in this particular implementation, there are no designated coordinator
	// nodes for a particular key (replicas are randomly chosen).
	//
	// This is set when this replicator is configured and can be swapped out for a function that
	// returns hard-coded nodes during testing.
	replicaChooser func(numReplicas int, exclude []uint64) ([]uint64, error)

	// Observability
	log *log.Logger

	// The public-facing API that this replicator must implement
	pb.ReplicatorServer

	// These functions are the internal, private RPCs for a node partaking in a basic leaderless
	// replication strategy
	pb.BasicLeaderlessReplicatorServer
}

// Args configures a given instance of leaderless replication.
type Args[T conflict.Clock] struct {
	Node *node.Node

	LocalStore store.Store[*conflict.KV[T]]

	ConflictResolver conflict.Resolver[T]

	// Quorum
	W int
	R int
}

// Configure is called by the orchestrator to start this node
//
// The "args" are any to support any replicator that might need arbitrary
// set of configuration values.
func Configure[T conflict.Clock](args any) *State[T] {
	a := args.(Args[T])

	s := &State[T]{
		node: a.Node,

		localStore:       a.LocalStore,
		conflictResolver: a.ConflictResolver,

		W: a.W,
		R: a.R,

		log: a.Node.Log,
	}
	s.replicaChooser = s.selectKRandomPeerNodeIDs

	// Configure lifecycle functions
	s.onReplicatorStart = func(node *node.Node) {
		s.conflictResolver.ReplicatorDidStart(node)
	}
	s.onMessageReceived = func(clock T) {
		s.conflictResolver.OnMessageReceive(clock)
	}
	s.onMessageSend = func() {
		s.conflictResolver.OnMessageSend()
	}
	s.onReplicatorStart(s.node)

	// gRPC public and private servers
	s.log.Printf("Starting gRPC server at %s", s.node.Addr.Host)
	grpcServer := a.Node.GrpcServer
	pb.RegisterReplicatorServer(grpcServer, s)
	pb.RegisterBasicLeaderlessReplicatorServer(grpcServer, s)
	go grpcServer.Serve(s.node.Listener)

	return s
}

// withRetries is a wrapper function to try the function f at most numRetries times, until f runs
// without returning an error. withRetries returns nil if f ran successfully, otherwise
// it returns the last error.
func (s *State[T]) withRetries(f func() error, numRetries int) error {
	var err error
	for i := 0; i < numRetries; i++ {
		err = f()
		if i > 0 {
			s.log.Printf("On retry %d, err %v", i, err)
		}

		if err == nil {
			return nil
		}

		s.log.Printf("Failed to establish connection: %v. %d retries left", err, numRetries-i-1)
	}
	return err
}

// dispatchToPeers chooses a set of random peers and calls function f for each of them. If a
// call to f fails, a new peer is chosen and the function f is retried on it.
//
// The first chosen peer is always the current (local) node.
func (s *State[T]) dispatchToPeers(ctx context.Context, num int, f func(ctx context.Context, replicaNodeID uint64) error) error {
	if num <= 0 {
		return errors.New("num must be positive")
	}

	alreadyChosenNodes := []uint64{s.node.ID}
	alreadyChosenNodesMu := sync.Mutex{}

	replicaNodeIDs, err := s.replicaChooser(num-1, alreadyChosenNodes)
	if err != nil {
		return err
	}
	alreadyChosenNodes = append(alreadyChosenNodes, replicaNodeIDs...)
	replicaNodeIDs = append(replicaNodeIDs, s.node.ID)

	s.log.Printf("dispatching to replicas %v", replicaNodeIDs)

	// true is sent on this channel whenever a goroutine to contact another node succeeds
	successC := make(chan bool)
	// an error is sent when a goroutine contacting another node faces an irrecoverable error
	errorC := make(chan error)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for _, replicaNodeID := range replicaNodeIDs {
		go func(replicaNodeID uint64) {
			for {
				err := f(ctx, replicaNodeID)
				if err != nil {
					// Try with a new replicaNodeID
					alreadyChosenNodesMu.Lock()

					newReplica, err := s.replicaChooser(1, alreadyChosenNodes)

					// If we cannot find another replica, we can't reach quorum for writing.
					// We have to give up.
					if err != nil {
						errorC <- err
						return
					}

					// Nobody else should be able to use newReplica[0] because we will try to
					// talk to it in the next loop iteration
					alreadyChosenNodes = append(alreadyChosenNodes, newReplica[0])
					alreadyChosenNodesMu.Unlock()

					replicaNodeID = newReplica[0]
				} else {
					successC <- true
					return
				}

				if ctx.Err() != nil {
					return
				}
			}
		}(replicaNodeID)
	}

	successCount := 0

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("communication timed out")
		case err = <-errorC:
			return err
		case <-successC:
			successCount++

			if successCount >= num {
				return nil
			}
		}
	}
}

// selectKRandomPeerNodeIDs selects k unique random peer nodes from s.PeerNodes, excluding nodes
// in the exclude slice.
func (s *State[T]) selectKRandomPeerNodeIDs(k int, exclude []uint64) ([]uint64, error) {
	peerNodes := s.node.PeerNodes
	n := len(peerNodes)

	isExcluded := func(ID uint64) bool {
		for _, excludedID := range exclude {
			if ID == excludedID {
				return true
			}
		}
		return false
	}

	if possibleNodes := n - len(exclude); k > possibleNodes {
		return nil, fmt.Errorf("cannot select %d peers if cluster with exclusions is of size %d", k, possibleNodes)
	}

	// Create a slice of all nodeIDs
	nodeIDs := make([]uint64, 0, n)
	for nodeID := range peerNodes {
		nodeIDs = append(nodeIDs, nodeID)
	}

	var peerIDs []uint64
	startingIndex := rand.Intn(n)

	// Iterate in order through the nodes, starting at a random node, until we
	for i := 0; i < n && len(peerIDs) < k; i++ {
		peerID := nodeIDs[(startingIndex+i)%n]

		if !isExcluded(peerID) {
			peerIDs = append(peerIDs, peerID)
		}
	}

	return peerIDs, nil
}

// ensureClock ensures that the provided clock is not nil. If it is, it is initialized to a zero
// clock. A zero clock is used as for new clients, any value is newer than what they have seen.
func (s *State[T]) ensureClock(clockPtr **pb.Clock) {
	clock := *clockPtr

	if clock == nil {
		*clockPtr = s.conflictResolver.ZeroClock().Proto()
	}
}
