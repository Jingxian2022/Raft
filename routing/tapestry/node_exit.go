/*
 *  Brown University, CS138, Spring 2023
 *
 *  Purpose: Defines functions for a node leaving the Tapestry mesh, and
 *  transferring its stored locations to a new node.
 */

package tapestry

import (
	"context"
	"errors"
	pb "modist/proto"
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
	local.NotifyLeave(context.Background(), &pb.LeaveNotification{
		From:        local.Id.String(),
		Replacement: "", // FIXME: Implement me!
	})
	local.blobstore.DeleteAll()
	go local.Node.GrpcServer.GracefulStop()
	return errors.New("Leave has not been implemented yet!")
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
	// Remove references to the `from` node from our routing table and backpointers
	local.Table.Remove(from)
	local.Backpointers.Remove(from)
	// Replacement can be an empty string so we don't want to parse it here
	replacement := leaveNotification.Replacement

	local.log.Printf(
		"Received leave notification from %v with replacement node %v\n",
		from,
		replacement,
	)

	// If replacement is not an empty string, add replacement to our routing table
	if replacement != "" {
		replacementID, err := ParseID(replacement)
		if err != nil {
			return nil, err
		}
		local.Table.Add(replacementID)
	}

	// TODO(students): [Tapestry] Implement me!
	return nil, errors.New("NotifyLeave has not been implemented yet!")
}
