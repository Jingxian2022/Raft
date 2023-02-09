/*
 *  Brown University, CS138, Spring 2023
 *
 *  Purpose: Implements functions that are invoked by other nodes over RPC.
 */

package tapestry

import (
	"errors"

	pb "modist/proto"

	"golang.org/x/net/context"
)

/**
 * RPC receiver functions
 */

func (local *Node) HelloCaller(ctx context.Context, n *pb.NodeMsg) (*pb.NodeMsg, error) {
	return local.Node.toNodeMsg(), nil
}

func (local *Node) FindRootCaller(ctx context.Context, id *pb.IdMsg) (*pb.RootMsg, error) {
	idVal, err := ParseID(id.Id)
	if err != nil {
		return nil, err
	}
	next, tr, err := local.FindRoot(idVal, id.Level)
	rsp := &pb.RootMsg{
		Next:     next.toNodeMsg(),
		ToRemove: remoteNodesToNodeMsgs(tr.Nodes()),
	}
	return rsp, err
}

func (local *Node) RegisterCaller(ctx context.Context, r *pb.Registration) (*pb.Ok, error) {
	// TODO(students): [Tapestry] Implement me!
	return nil, nil
}

func (local *Node) FetchCaller(ctx context.Context, key *pb.ReplicaGroupID) (*pb.FetchedLocations, error) {
	isRoot, values := local.Fetch(key.Key)

	rsp := &pb.FetchedLocations{
		Values: remoteNodesToNodeMsgs(values),
		IsRoot: isRoot,
	}
	return rsp, nil
}

func (local *Node) RemoveBadNodesCaller(ctx context.Context, nodes *pb.Neighbors) (*pb.Ok, error) {
	// TODO(students): [Tapestry] Implement me!
	return nil, nil
}

func (local *Node) AddNodeCaller(ctx context.Context, n *pb.NodeMsg) (*pb.Neighbors, error) {
	neighbors, err := local.AddNode(toRemoteNode(n))

	rsp := &pb.Neighbors{
		Neighbors: remoteNodesToNodeMsgs(neighbors),
	}
	return rsp, err
}

func (local *Node) AddNodeMulticastCaller(ctx context.Context, m *pb.MulticastRequest) (*pb.Neighbors, error) {
	// TODO(students): [Tapestry] Implement me!
	return nil, nil
}

func (local *Node) TransferCaller(ctx context.Context, td *pb.TransferData) (*pb.Ok, error) {
	parsedData := make(map[string][]RemoteNode)
	for key, set := range td.Data {
		parsedData[key] = nodeMsgsToRemoteNodes(set.Neighbors)
	}
	err := local.Transfer(toRemoteNode(td.From), parsedData)

	rsp := &pb.Ok{
		Ok: true,
	}
	return rsp, err
}

func (local *Node) AddBackpointerCaller(ctx context.Context, n *pb.NodeMsg) (*pb.Ok, error) {
	// TODO(students): [Tapestry] Implement me!
	return nil, nil
}

func (local *Node) RemoveBackpointerCaller(ctx context.Context, n *pb.NodeMsg) (*pb.Ok, error) {
	err := local.RemoveBackpointer(toRemoteNode(n))
	rsp := &pb.Ok{
		Ok: true,
	}
	return rsp, err
}

func (local *Node) GetBackpointersCaller(ctx context.Context, br *pb.BackpointerRequest) (*pb.Neighbors, error) {
	// TODO(students): [Tapestry] Implement me!
	return nil, nil
}

func (local *Node) NotifyLeaveCaller(ctx context.Context, ln *pb.LeaveNotification) (*pb.Ok, error) {
	replacement := toRemoteNode(ln.Replacement)
	err := local.NotifyLeave(toRemoteNode(ln.From), &replacement)
	rsp := &pb.Ok{
		Ok: true,
	}
	return rsp, err
}

func (local *Node) BlobStoreFetchCaller(ctx context.Context, key *pb.ReplicaGroupID) (*pb.DataBlob, error) {
	data, isOk := local.blobstore.Get(key.Key)
	var err error
	if !isOk {
		err = errors.New("Key not found")
	}
	return &pb.DataBlob{
		Key:  key.Key,
		Data: data,
	}, err
}

func (local *Node) TapestryLookupCaller(ctx context.Context, key *pb.ReplicaGroupID) (*pb.Neighbors, error) {
	nodes, err := local.Lookup(key.Key)
	return &pb.Neighbors{
		Neighbors: remoteNodesToNodeMsgs(nodes),
	}, err
}

func (local *Node) TapestryStoreCaller(ctx context.Context, blob *pb.DataBlob) (*pb.Ok, error) {
	return &pb.Ok{Ok: true}, local.Store(blob.Key, blob.Data)
}

func remoteNodesToNodeMsgs(remoteNodes []RemoteNode) []*pb.NodeMsg {
	nodeMsgs := make([]*pb.NodeMsg, len(remoteNodes))
	for i, thing := range remoteNodes {
		nodeMsgs[i] = thing.toNodeMsg()
	}
	return nodeMsgs
}

func nodeMsgsToRemoteNodes(nodeMsgs []*pb.NodeMsg) []RemoteNode {
	remoteNodes := make([]RemoteNode, len(nodeMsgs))
	for i, thing := range nodeMsgs {
		remoteNodes[i] = toRemoteNode(thing)
	}
	return remoteNodes
}
