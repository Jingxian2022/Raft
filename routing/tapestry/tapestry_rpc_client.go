/*
 *  Brown University, CS138, Spring 2022
 *
 *  Purpose: Provides wrappers around the client interface of GRPC to invoke
 *  functions on remote tapestry nodes.
 */

package tapestry

import (
	"sync"
	"time"

	// util "github.com/brown-csci1380/tracing-framework-go/xtrace/grpcutil"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "modist/proto"
)

const GRPCTimeout = 5 * time.Second

// clientUnaryInterceptor is a client unary interceptor that injects a default timeout
func clientUnaryInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	ctx, cancel := context.WithTimeout(ctx, GRPCTimeout)
	defer cancel()

	return invoker(ctx, method, req, reply, cc, opts...)
}

// RemoteNode represents non-local node addresses in the tapestry
type RemoteNode struct {
	ID      ID
	Address string
}

// Turns a NodeMsg into a RemoteNode
func toRemoteNode(n *pb.NodeMsg) RemoteNode {
	if n == nil {
		return RemoteNode{}
	}
	idVal, err := ParseID(n.Id)
	if err != nil {
		return RemoteNode{}
	}
	return RemoteNode{
		ID:      idVal,
		Address: n.Address,
	}
}

// Turns a RemoteNode into a NodeMsg
func (n *RemoteNode) toNodeMsg() *pb.NodeMsg {
	if n == nil {
		return nil
	}
	return &pb.NodeMsg{
		Id:      n.ID.String(),
		Address: n.Address,
	}
}

/**
 *  RPC invocation functions
 */

var connMap = make(map[string]*grpc.ClientConn)
var connMapLock = &sync.RWMutex{}

func CloseAllConnections() {
	connMapLock.Lock()
	defer connMapLock.Unlock()
	for k, conn := range connMap {
		conn.Close()
		delete(connMap, k)
	}
}

// Creates a new client connection to the given remote node
func makeClientConn(remote *RemoteNode) (*grpc.ClientConn, error) {
	dialOptions := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithUnaryInterceptor(clientUnaryInterceptor)}
	return grpc.Dial(remote.Address, dialOptions...)
}

// Creates or returns a cached RPC client for the given remote node
func (remote *RemoteNode) ClientConn() (pb.TapestryRPCClient, error) {
	connMapLock.RLock()
	if cc, ok := connMap[remote.Address]; ok {
		connMapLock.RUnlock()
		return pb.NewTapestryRPCClient(cc), nil
	}
	connMapLock.RUnlock()

	cc, err := makeClientConn(remote)
	if err != nil {
		return nil, err
	}
	connMapLock.Lock()
	connMap[remote.Address] = cc
	connMapLock.Unlock()

	return pb.NewTapestryRPCClient(cc), err
}

// Remove the client connection to the given node, if present
func (remote *RemoteNode) RemoveClientConn() {
	connMapLock.Lock()
	defer connMapLock.Unlock()
	if cc, ok := connMap[remote.Address]; ok {
		cc.Close()
		delete(connMap, remote.Address)
	}
}

// Check the error and remove the client connection if necessary
func (remote *RemoteNode) connCheck(err error) error {
	if err != nil {
		remote.RemoveClientConn()
	}
	return err
}

// Say hello to a remote address, and get the tapestry node there
func SayHelloRPC(addr string, joiner RemoteNode) (RemoteNode, error) {
	remote := &RemoteNode{Address: addr}
	cc, err := remote.ClientConn()
	if err != nil {
		return RemoteNode{}, err
	}
	node, err := cc.HelloCaller(context.Background(), joiner.toNodeMsg())
	return toRemoteNode(node), remote.connCheck(err)
}

func (remote *RemoteNode) FindRootRPC(id ID, level int32) (RemoteNode, *NodeSet, error) {
	// TODO(students): [Tapestry] Implement me!
	return RemoteNode{}, nil, nil
}

func (remote *RemoteNode) RegisterRPC(key string, replica RemoteNode) (bool, error) {
	cc, err := remote.ClientConn()
	if err != nil {
		return false, err
	}
	rsp, err := cc.RegisterCaller(context.Background(), &pb.Registration{
		FromNode: replica.toNodeMsg(),
		Key:      key,
	})
	return rsp.GetOk(), remote.connCheck(err)
}

func (remote *RemoteNode) FetchRPC(key string) (bool, []RemoteNode, error) {
	// TODO(students): [Tapestry] Implement me!
	return false, nil, nil
}

func (remote *RemoteNode) RemoveBadNodesRPC(badnodes []RemoteNode) error {
	cc, err := remote.ClientConn()
	if err != nil {
		return err
	}
	_, err = cc.RemoveBadNodesCaller(context.Background(), &pb.Neighbors{Neighbors: remoteNodesToNodeMsgs(badnodes)})
	return remote.connCheck(err)
}

func (remote *RemoteNode) AddNodeRPC(toAdd RemoteNode) ([]RemoteNode, error) {
	// TODO(students): [Tapestry] Implement me!
	return nil, nil
}

func (remote *RemoteNode) AddNodeMulticastRPC(newNode RemoteNode, level int) ([]RemoteNode, error) {
	cc, err := remote.ClientConn()
	if err != nil {
		return nil, err
	}
	rsp, err := cc.AddNodeMulticastCaller(context.Background(), &pb.MulticastRequest{
		NewNode: newNode.toNodeMsg(),
		Level:   int32(level),
	})
	if err != nil {
		return nil, remote.connCheck(err)
	}
	return nodeMsgsToRemoteNodes(rsp.Neighbors), remote.connCheck(err)
}

func (remote *RemoteNode) TransferRPC(from RemoteNode, data map[string][]RemoteNode) error {
	// TODO(students): [Tapestry] Implement me!
	return nil
}

func (remote *RemoteNode) AddBackpointerRPC(bp RemoteNode) error {
	cc, err := remote.ClientConn()
	if err != nil {
		return err
	}
	_, err = cc.AddBackpointerCaller(context.Background(), bp.toNodeMsg())
	return remote.connCheck(err)
}

func (remote *RemoteNode) RemoveBackpointerRPC(bp RemoteNode) error {
	// TODO(students): [Tapestry] Implement me!
	return nil
}

func (remote *RemoteNode) GetBackpointersRPC(from RemoteNode, level int) ([]RemoteNode, error) {
	cc, err := remote.ClientConn()
	if err != nil {
		return nil, err
	}
	rsp, err := cc.GetBackpointersCaller(context.Background(), &pb.BackpointerRequest{
		From:  from.toNodeMsg(),
		Level: int32(level),
	})
	if err != nil {
		return nil, remote.connCheck(err)
	}
	return nodeMsgsToRemoteNodes(rsp.Neighbors), remote.connCheck(err)
}

func (remote *RemoteNode) NotifyLeaveRPC(from RemoteNode, replacement *RemoteNode) error {
	// TODO: (Tapestry) students should implement this
	return nil
}

func (remote *RemoteNode) BlobStoreFetchRPC(key string) (*[]byte, error) {
	cc, err := remote.ClientConn()
	if err != nil {
		return nil, err
	}
	rsp, err := cc.BlobStoreFetchCaller(context.Background(), &pb.ReplicaGroupID{Key: key})
	if err != nil {
		return nil, remote.connCheck(err)
	}
	return &rsp.Data, remote.connCheck(err)
}

func (remote *RemoteNode) TapestryLookupRPC(key string) ([]RemoteNode, error) {
	cc, err := remote.ClientConn()
	if err != nil {
		return nil, err
	}
	rsp, err := cc.TapestryLookupCaller(context.Background(), &pb.ReplicaGroupID{Key: key})
	if err != nil {
		return nil, remote.connCheck(err)
	}
	return nodeMsgsToRemoteNodes(rsp.Neighbors), remote.connCheck(err)
}

func (remote *RemoteNode) TapestryStoreRPC(key string, value []byte) error {
	cc, err := remote.ClientConn()
	if err != nil {
		return err
	}
	_, err = cc.TapestryStoreCaller(context.Background(), &pb.DataBlob{
		Key:  key,
		Data: value,
	})
	return remote.connCheck(err)
}
