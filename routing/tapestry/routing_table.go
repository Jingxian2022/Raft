/*
 *  Brown University, CS138, Spring 2023
 *
 *  Purpose: Defines the RoutingTable type and provides methods for interacting
 *  with it.
 */

package tapestry

import (
	"fmt"
	"sync"
)

// RoutingTable has a number of levels equal to the number of digits in an ID
// (default 40). Each level has a number of slots equal to the digit base
// (default 16). A node that exists on level n thereby shares a prefix of length
// n with the local node. Access to the routing table is protected by a mutex.
type RoutingTable struct {
	localId ID                 // The ID of the local tapestry node
	Rows    [DIGITS][BASE][]ID // The rows of the routing table (stores IDs of remote tapestry nodes)
	mutex   sync.Mutex         // To manage concurrent access to the routing table (could also have a per-level mutex)
}

// NewRoutingTable creates and returns a new routing table, placing the local node at the
// appropriate slot in each level of the table.
func NewRoutingTable(me ID) *RoutingTable {
	t := new(RoutingTable)
	t.localId = me

	// Create the node lists with capacity of SLOTSIZE
	for i := 0; i < DIGITS; i++ {
		for j := 0; j < BASE; j++ {
			t.Rows[i][j] = make([]ID, 0, SLOTSIZE)
		}
	}

	// Make sure each row has at least our node in it
	for i := 0; i < DIGITS; i++ {
		slot := t.Rows[i][t.localId[i]]
		t.Rows[i][t.localId[i]] = append(slot, t.localId)
	}

	return t
}

// Add adds the given node to the routing table.
//
// Note you should not add the node to preceding levels. You need to add the node
// to one specific slot in the routing table (or replace an element if the slot is full
// at SLOTSIZE).
//
// Returns true if the node did not previously exist in the table and was subsequently added.
// Returns the previous node in the table, if one was overwritten.
func (t *RoutingTable) Add(remoteNodeId ID) (added bool, previous *ID) {

	t.mutex.Lock()
	defer t.mutex.Unlock()

	// TODO(students): [Tapestry] Implement me!
	n := SharedPrefixLength(t.localId, remoteNodeId)
	slot := t.Rows[n][remoteNodeId[n]]
	slotLength := len(slot)

	fmt.Printf("Adding to slot %v on row %v\n", remoteNodeId[n], n)

	pos := slotLength
	for k := slotLength - 1; k >= 0; k-- {
		// if already exists
		if slot[k] == remoteNodeId {
			return false, nil
		}
		if t.localId.Closer(remoteNodeId, slot[k]) {
			// replace last one with new node when slot is full
			if k == SLOTSIZE-1 {
				previous = &slot[k]
			} else if k == slotLength-1 { // copy slot[k] to slot[k+1]
				slot = append(slot, slot[k])
			} else {
				slot[k+1] = slot[k]
			}
			pos = k // pos to insert new node ID
		} else {
			break
		}
	}

	if pos == slotLength && pos != SLOTSIZE {
		slot = append(slot, remoteNodeId)
		added = true
		fmt.Printf("case 1\n")
	} else if pos != SLOTSIZE {
		slot[pos] = remoteNodeId
		added = true
		fmt.Printf("case 2\n")
	}
	t.Rows[n][remoteNodeId[n]] = slot
	return added, previous
}

// Remove removes the specified node from the routing table, if it exists.
// Returns true if the node was in the table and was successfully removed.
// Return false if a node tries to remove itself from the table.
func (t *RoutingTable) Remove(remoteNodeId ID) (wasRemoved bool) {

	t.mutex.Lock()
	defer t.mutex.Unlock()

	// TODO(students): [Tapestry] Implement me!
	if t.localId == remoteNodeId {
		return false
	}
	n := SharedPrefixLength(t.localId, remoteNodeId)
	slot := t.Rows[n][remoteNodeId[n]]
	pos := -1
	for k := 0; k < len(slot); k++ {
		if slot[k] == remoteNodeId {
			pos = k
			break
		}
	}
	if pos != -1 {
		slot = append(slot[:pos], slot[pos+1:]...)
		wasRemoved = true
	}
	t.Rows[n][remoteNodeId[n]] = slot
	return wasRemoved
}

// GetLevel gets ALL nodes on the specified level of the routing table, EXCLUDING the local node.
func (t *RoutingTable) GetLevel(level int) (nodeIds []ID) {

	t.mutex.Lock()
	defer t.mutex.Unlock()

	// TODO(students): [Tapestry] Implement me!

	nodeIds = make([]ID, 0)
	for k := 0; k < BASE; k++ {
		for i := 0; i < len(t.Rows[level][k]); i++ {
			if t.Rows[level][k][i] != t.localId {
				nodeIds = append(nodeIds, t.Rows[level][k][i])
			}
		}
	}

	return nodeIds
}

// FindNextHop searches the table for the closest next-hop node for the provided ID starting at the given level.
func (t *RoutingTable) FindNextHop(id ID, level int32) ID {

	t.mutex.Lock()
	defer t.mutex.Unlock()

	// TODO(students): [Tapestry] Implement me!
	for currentLevel := level; currentLevel < DIGITS; currentLevel++ {
		idx := id[currentLevel]
		for k := 0; k < BASE; k++ {
			// Finds a non-empty cell in the table
			if len(t.Rows[currentLevel][idx]) != 0 {
				// If the first node in the non-empty cell is non-local, return it.
				if t.Rows[currentLevel][idx][0] != t.localId {
					return t.Rows[currentLevel][idx][0]
				}
				// Jump down to the next level
				break
			}
			// Move right along the level starting from that slot
			idx = (idx + 1) % BASE
		}
	}

	// for currentLevel := level; currentLevel < DIGITS; currentLevel++ {
	// 	flag_jumplevel := false
	// 	// iterate through the slots in the level
	// 	for currentSlot := id[currentLevel]; currentSlot < BASE; currentSlot++ {
	// 		// Finds a non-empty cell in the table
	// 		if len(t.Rows[currentLevel][currentSlot]) != 0 {
	// 			// iterate through the nodes in the cell
	// 			for i := 0; i < len(t.Rows[currentLevel][currentSlot]); i++ {
	// 				// If the first node in the non-empty cell is non-local, return it.
	// 				if t.Rows[currentLevel][currentSlot][i] != t.localId {
	// 					return t.Rows[currentLevel][currentSlot][i]
	// 				}
	// 				//  If the first node in the non-empty cell is local, we need to keep searching.
	// 				// Jump down to the first slot of the next level
	// 				flag_jumplevel = true
	// 				break // jump to next level
	// 			}
	// 		}
	// 		if flag_jumplevel {
	// 			break // jump to next level
	// 		}
	// 	}
	// 	if flag_jumplevel {
	// 		continue
	// 	}
	// 	return t.localId
	// }
	return t.localId
}
