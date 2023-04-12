/*
 *  Brown University, CS138, Spring 2023
 *
 *  Purpose: Defines global constants and functions to create and join a new
 *  node into a Tapestry mesh, and functions for altering the routing table
 *  and backpointers of the local node that are invoked over RPC.
 */

package tapestry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"modist/orchestrator/node"
	pb "modist/proto"
	"sort"
	"sync"
	"time"
)

// BASE is the base of a digit of an ID.  By default, a digit is base-16.
const BASE = 16

// DIGITS is the number of digits in an ID.  By default, an ID has 40 digits.
const DIGITS = 40

// RETRIES is the number of retries on failure. By default we have 3 retries.
const RETRIES = 3

// K is neigborset size during neighbor traversal before fetching backpointers. By default this has a value of 10.
const K = 10

// SLOTSIZE is the size each slot in the routing table should store this many nodes. By default this is 3.
const SLOTSIZE = 3

// REPUBLISH is object republish interval for nodes advertising objects.
const REPUBLISH = 10 * time.Second

// TIMEOUT is object timeout interval for nodes storing objects.
const TIMEOUT = 25 * time.Second

// TapestryNode is the main struct for the local Tapestry node. Methods can be invoked locally on this struct.
type TapestryNode struct {
	Node *node.Node // Node that this Tapestry node is part of
	Id   ID         // ID of node in the form of a slice (makes it easier to iterate)

	Table          *RoutingTable // The routing table
	Backpointers   *Backpointers // Backpointers to keep track of other nodes that point to us
	LocationsByKey *LocationMap  // Stores keys for which this node is the root
	blobstore      *BlobStore    // Stores blobs on the local node

	// Observability
	log *log.Logger

	// These functions are the internal, private RPCs for a routing node using Tapestry
	pb.UnsafeTapestryRPCServer
}

func (local *TapestryNode) String() string {
	return fmt.Sprint(local.Id)
}

// Called in tapestry initialization to create a tapestry node struct
func newTapestryNode(node *node.Node) *TapestryNode {
	tn := new(TapestryNode)

	tn.Node = node
	tn.Id = MakeID(node.ID)
	tn.Table = NewRoutingTable(tn.Id)
	tn.Backpointers = NewBackpointers(tn.Id)
	tn.LocationsByKey = NewLocationMap()
	tn.blobstore = NewBlobStore()

	tn.log = node.Log

	return tn
}

// Start Tapestry Node
func StartTapestryNode(node *node.Node, connectTo uint64, join bool) (tn *TapestryNode, err error) {
	// Create the local node
	tn = newTapestryNode(node)

	tn.log.Printf("Created tapestry node %v\n", tn)

	grpcServer := tn.Node.GrpcServer
	pb.RegisterTapestryRPCServer(grpcServer, tn)

	// If specified, connect to the provided ID
	if join {
		// If provided ID doesn't exist, return an error
		if _, ok := node.PeerConns[connectTo]; !ok {
			return nil, fmt.Errorf(
				"Error joining Tapestry node with id %v; Unable to find node %v in peerConns",
				connectTo,
				connectTo,
			)
		}

		err = tn.Join(MakeID(connectTo))
		if err != nil {
			tn.log.Printf(err.Error())
			return nil, err
		}
	}

	return tn, nil
}

// Join is invoked when starting the local node, if we are connecting to an existing Tapestry.
//
// - Find the root for our node's ID
// - Call AddNode on our root to initiate the multicast and receive our initial neighbor set. Add them to our table.
// - Iteratively get backpointers from the neighbor set for all levels in range [0, SharedPrefixLength]
// and populate routing table
func (local *TapestryNode) Join(remoteNodeId ID) error {
	local.log.Println("Joining tapestry node", remoteNodeId)

	// Route to our root
	rootIdPtr, err := local.FindRootOnRemoteNode(remoteNodeId, local.Id)
	if err != nil {
		return fmt.Errorf("Error joining existing tapestry node %v, reason: %v", remoteNodeId, err)
	}
	rootId := *rootIdPtr

	// Add ourselves to our root by invoking AddNode on the remote node
	nodeMsg := &pb.NodeMsg{
		Id: local.Id.String(),
	}

	conn := local.Node.PeerConns[local.RetrieveID(rootId)]
	rootNode := pb.NewTapestryRPCClient(conn)
	resp, err := rootNode.AddNode(context.Background(), nodeMsg)
	if err != nil {
		return fmt.Errorf("Error adding ourselves to root node %v, reason: %v", rootId, err)
	}

	// Add the neighbors to our local routing table.
	neighborIds, err := stringSliceToIds(resp.Neighbors)
	if err != nil {
		return fmt.Errorf("Error parsing neighbor IDs, reason: %v", err)
	}

	for _, neighborId := range neighborIds {
		local.AddRoute(neighborId)
	}

	// TODO(students): [Tapestry] Implement me!
	n := SharedPrefixLength(local.Id, remoteNodeId)
	err = local.traverseBackpointers(neighborIds, n)
	return err
}

func (local *TapestryNode) traverseBackpointers(neighbors []ID, level int) error {
	for k := level; k >= 0; k-- {
		sort.SliceStable(neighbors, func(i, j int) bool {
			return local.Id.Closer(neighbors[i], neighbors[j])
		})
		if len(neighbors) > K {
			neighbors = neighbors[:K]
		}
		nextNeighbors := idsToStringSlice(neighbors)
		for _, neighbor := range neighbors {
			go func(neighbor ID, nextNeighbors *[]string) {
				conn := local.Node.PeerConns[local.RetrieveID(neighbor)]
				neighborNode := pb.NewTapestryRPCClient(conn)
				backPointers, err := neighborNode.GetBackpointers(context.Background(), &pb.BackpointerRequest{From: local.Id.String(), Level: int32(k)})
				if err != nil {
					return
				}
				*nextNeighbors = append(*nextNeighbors, backPointers.Neighbors...)
			}(neighbor, &nextNeighbors)
		}
		nextNeighbors = Merge(nextNeighbors)
		for _, neighbor := range nextNeighbors {
			neighborId, err := ParseID(neighbor)
			if err != nil {
				return errors.New("Failed to add to routing table!")
			}
			local.AddRoute(neighborId)
		}
		neighbors, _ = stringSliceToIds(nextNeighbors)
	}
	return nil
}

// AddNode adds node to the tapestry
//
// - Begin the acknowledged multicast
// - Return the neighborset from the multicast
func (local *TapestryNode) AddNode(
	ctx context.Context,
	nodeMsg *pb.NodeMsg,
) (*pb.Neighbors, error) {
	nodeId, err := ParseID(nodeMsg.Id)
	if err != nil {
		return nil, err
	}

	local.log.Printf("New node %v, old node %v\n", nodeId, local.Id)

	multicastRequest := &pb.MulticastRequest{
		NewNode: nodeMsg.Id,
		Level:   int32(SharedPrefixLength(nodeId, local.Id)),
	}
	return local.AddNodeMulticast(context.Background(), multicastRequest)
}

func (local *TapestryNode) Output() {
	for i := 0; i < DIGITS; i++ {
		for j := 0; j < BASE; j++ {
			for k := 0; k < len(local.Table.Rows[i][j]); k++ {
				local.log.Printf("Row %v slot %v element %v: %v\n", i, j, k, local.Table.Rows[i][j][k])
			}
		}
	}
}

// AddNodeMulticast sends newNode to need-to-know nodes participating in the multicast.
//   - Perform multicast to need-to-know nodes
//   - Add the route for the new node (use `local.addRoute`)
//   - Transfer of appropriate router info to the new node (use `local.locationsByKey.GetTransferRegistrations`)
//     If error, rollback the location map (add back unsuccessfully transferred objects)
//
// - Propagate the multicast to the specified row in our routing table and await multicast responses
// - Return the merged neighbor set
//
// - note: `local.table.GetLevel` does not return the local node, so you must manually add this to the neighbors set
func (local *TapestryNode) AddNodeMulticast(
	ctx context.Context,
	multicastRequest *pb.MulticastRequest,
) (*pb.Neighbors, error) {
	newNodeId, err := ParseID(multicastRequest.NewNode)
	if err != nil {
		return nil, err
	}
	level := int(multicastRequest.Level)

	local.log.Printf("Add node multicast at node %v of new node %v at level %v\n", local.Id, newNodeId, level)

	// TODO(students): [Tapestry] Implement me!
	targets := local.Table.GetLevel(level)
	targets = append(targets, local.Id)
	results := make([]string, 0)

	mtx := sync.Mutex{}
	wg := sync.WaitGroup{}
	if level < DIGITS {
		for _, target := range targets {
			wg.Add(1)
			go func(target ID, results *[]string) {
				defer wg.Done()
				local.log.Printf("To call add node multicast at node %v of new node %v at level %v\n", target.String(), newNodeId, level+1)

				conn := local.Node.PeerConns[local.RetrieveID(target)]
				targetNode := pb.NewTapestryRPCClient(conn)
				targetNeighbors, _ := targetNode.AddNodeMulticast(context.Background(), &pb.MulticastRequest{NewNode: multicastRequest.GetNewNode(), Level: multicastRequest.Level + 1})
				mtx.Lock()
				*results = append(*results, targetNeighbors.Neighbors...)
				mtx.Unlock()
			}(target, &results)

		}
		//local.log.Printf("wait1")
		wg.Wait()
		//local.log.Printf("wait2")

		results = Merge(results)
		for _, res := range results {
			if res == local.Id.String() {
				return &pb.Neighbors{Neighbors: results}, nil
			}
		}

		//local.Output()
		err = local.AddRoute(newNodeId)
		if err != nil {
			return nil, err
		}
		go func() {
			objects := local.LocationsByKey.GetTransferRegistrations(local.Id, newNodeId)
			conn := local.Node.PeerConns[local.RetrieveID(newNodeId)]
			newNode := pb.NewTapestryRPCClient(conn)

			transferData := pb.TransferData{From: local.Id.String()}
			for key, set := range objects {
				transferData.Data[key] = &pb.Neighbors{Neighbors: idsToStringSlice(set)}
			}
			_, err = newNode.Transfer(ctx, &transferData)
			if err != nil {
				local.LocationsByKey.RegisterAll(objects, TIMEOUT)
			}
		}()

	}

	neighbors := pb.Neighbors{Neighbors: results}
	local.log.Printf("Add node multicast %v at level %v DONE\n", newNodeId, level)
	return &neighbors, nil
}

func Merge(input []string) (output []string) {
	tmp := map[string]struct{}{}
	for _, str := range input {
		tmp[str] = struct{}{}
	}
	for str := range tmp {
		output = append(output, str)
	}
	return output
}

// AddBackpointer adds the from node to our backpointers, and possibly add the node to our
// routing table, if appropriate
func (local *TapestryNode) AddBackpointer(
	ctx context.Context,
	nodeMsg *pb.NodeMsg,
) (*pb.Ok, error) {
	local.log.Printf("Error 2")

	id, err := ParseID(nodeMsg.Id)
	if err != nil {
		return nil, err
	}
	local.log.Printf("Error 1")

	if local.Backpointers.Add(id) {
		local.log.Printf("Added backpointer %v\n", id)
	}
	local.AddRoute(id)

	ok := &pb.Ok{
		Ok: true,
	}
	return ok, nil
}

// RemoveBackpointer removes the from node from our backpointers
func (local *TapestryNode) RemoveBackpointer(
	ctx context.Context,
	nodeMsg *pb.NodeMsg,
) (*pb.Ok, error) {
	id, err := ParseID(nodeMsg.Id)
	if err != nil {
		return nil, err
	}

	if local.Backpointers.Remove(id) {
		local.log.Printf("Removed backpointer %v\n", id)
	}

	ok := &pb.Ok{
		Ok: true,
	}
	return ok, nil
}

// GetBackpointers gets all backpointers at the level specified, and possibly adds the node to our
// routing table, if appropriate
func (local *TapestryNode) GetBackpointers(
	ctx context.Context,
	backpointerReq *pb.BackpointerRequest,
) (*pb.Neighbors, error) {
	id, err := ParseID(backpointerReq.From)
	if err != nil {
		return nil, err
	}
	level := int(backpointerReq.Level)

	local.log.Printf("Sending level %v backpointers to %v\n", level, id)
	backpointers := local.Backpointers.Get(level)
	err = local.AddRoute(id)
	if err != nil {
		return nil, err
	}

	resp := &pb.Neighbors{
		Neighbors: idsToStringSlice(backpointers),
	}
	return resp, err
}

// RemoveBadNodes discards all the provided nodes
// - Remove each node from our routing table
// - Remove each node from our set of backpointers
func (local *TapestryNode) RemoveBadNodes(
	ctx context.Context,
	neighbors *pb.Neighbors,
) (*pb.Ok, error) {
	badnodes, err := stringSliceToIds(neighbors.Neighbors)
	if err != nil {
		return nil, err
	}

	for _, badnode := range badnodes {
		if local.Table.Remove(badnode) {
			local.log.Printf("Removed bad node %v\n", badnode)
		}
		if local.Backpointers.Remove(badnode) {
			local.log.Printf("Removed bad node backpointer %v\n", badnode)
		}
	}

	resp := &pb.Ok{
		Ok: true,
	}
	return resp, nil
}

// Utility function that adds a node to our routing table.
//
// - Adds the provided node to the routing table, if appropriate.
// - If the node was added to the routing table, notify the node of a backpointer
// - If an old node was removed from the routing table, notify the old node of a removed backpointer
func (local *TapestryNode) AddRoute(remoteNodeId ID) error {
	// TODO(students): [Tapestry] Implement me!
	added, previous := local.Table.Add(remoteNodeId)
	//local.Output()
	if added {
		go func() {
			conn := local.Node.PeerConns[local.RetrieveID(remoteNodeId)]
			remoteNode := pb.NewTapestryRPCClient(conn)
			local.log.Printf("RemoteNodeID %v, local ID %v\n", remoteNodeId, local.Id.String())
			local.log.Printf("RemoteNode %v", remoteNode)
			_, err := remoteNode.AddBackpointer(context.Background(), &pb.NodeMsg{Id: local.Id.String()})
			if err != nil {
				panic(err)
				//return fmt.Errorf("Error adding backpointer.")
			}
		}()
	}
	if previous != nil {
		go func() {
			conn := local.Node.PeerConns[local.RetrieveID(*previous)]
			previousNode := pb.NewTapestryRPCClient(conn)
			_, err := previousNode.RemoveBackpointer(context.Background(), &pb.NodeMsg{Id: local.String()})
			if err != nil {
				panic(err)
				//return fmt.Errorf("Error removing backpointer.")
			}
		}()
	}
	return nil
}
