/*
 *  Brown University, CS138, Spring 2023
 *
 *  Purpose: Defines functions for a node leaving the Tapestry mesh, and
 *  transferring its stored locations to a new node.
 */

package tapestry

import (
	"context"
	pb "modist/proto"
	"sync"
)

// Kill this node without gracefully leaving the tapestry.
func (local *TapestryNode) Kill() {
	local.blobstore.DeleteAll()
	local.Node.GrpcServer.Stop()
}

// Leave gracefully exits the Tapestry mesh.
//
// - Notify the nodes in our backpointers that we are leaving by calling NotifyLeave
// - If possible, give each backpointer a suitable alternative node from our routing table
func (local *TapestryNode) Leave() error {
	// TODO(students): [Tapestry] Implement me!
	// find the nodes in our backpointers
	replacement := ""
	for i := DIGITS - 1; i >= 0; i-- {
		replacement = ""
		if i != DIGITS-1 {
			nodes := local.Table.GetLevel(i)
			if len(nodes) > 0 {
				replacement = nodes[0].String()
			}
		}

		backPointers := local.Backpointers.Get(i)
		var wg sync.WaitGroup
		for _, backPointer := range backPointers {
			wg.Add(1)
			go func(remoteNodeId ID, replacement string) {
				defer wg.Done()
				conn := local.Node.PeerConns[local.RetrieveID(remoteNodeId)]
				toNotify := pb.NewTapestryRPCClient(conn)
				_, err := toNotify.NotifyLeave(context.Background(), &pb.LeaveNotification{
					From:        local.Id.String(),
					Replacement: replacement,
				})

				if err != nil {
					local.RemoveBadNodes(context.Background(), &pb.Neighbors{
						Neighbors: idsToStringSlice([]ID{remoteNodeId}),
					})
				}
			}(backPointer, replacement)
		}
		wg.Wait()
	}

	local.blobstore.DeleteAll()
	go local.Node.GrpcServer.GracefulStop()
	return nil
}

// NotifyLeave occurs when another node is informing us of a graceful exit.
// - Remove references to the `from` node from our routing table and backpointers
// - If replacement is not an empty string, add replacement to our routing table
func (local *TapestryNode) NotifyLeave(
	ctx context.Context,
	leaveNotification *pb.LeaveNotification,
) (*pb.Ok, error) {
	from, err := ParseID(leaveNotification.From)
	if err != nil {
		return nil, err
	}

	// Replacement can be an empty string so we don't want to parse it here
	replacement := leaveNotification.Replacement

	local.log.Printf(
		"Received leave notification from %v with replacement node %v\n",
		from,
		replacement,
	)

	// Remove references to the `from` node from our routing table and backpointers
	local.Table.Remove(from)
	local.Backpointers.Remove(from)

	// If replacement is not an empty string, add replacement to our routing table
	if replacement != "" {
		replacementID, err := ParseID(replacement)
		if err != nil {
			return nil, err
		}
		err = local.AddRoute(replacementID)
		if err != nil {
			return nil, err
		}
	}

	// TODO(students): [Tapestry] Implement me!
	return &pb.Ok{Ok: true}, nil
}
