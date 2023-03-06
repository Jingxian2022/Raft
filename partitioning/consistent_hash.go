package partitioning

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"

	"golang.org/x/exp/slices"
)

// Lookup returns the ID of the replica group to which the specified key is assigned.
// It also returns a hashed version of the key as the second return value.
//
// The replica group ID corresponding to a given key is found by looking up the first
// virtual node that succeeds the hashed key on the ring and returning the replica group ID
// to which this virtual node corresponds. If no replica groups have been added to the ring,
// an error is returned.
func (c *ConsistentHash) Lookup(key string) (id uint64, rewrittenKey string, err error) {

	// TODO(students): [Partitioning] Implement me!
	if len(c.virtualNodes) == 0 {
		return 0, "", errors.New("No implements")
	}
	hash := c.keyHash(key)
	mykey := virtualNode{hash: hash}

	if bytes.Compare(hash[:], c.virtualNodes[len(c.virtualNodes)-1].hash[:]) > 0 {
		return c.virtualNodes[0].id, hashToString(hash), nil
	}

	for i, node := range c.virtualNodes {
		if i == len(c.virtualNodes)-1 {
			if virtualNodeCmp(mykey, node) == 0 {
				id = node.id
			} else {
				id = c.virtualNodes[0].id
			}
			continue
		}
		if virtualNodeLess(mykey, node) || virtualNodeCmp(mykey, node) == 0 {
			id = node.id
			break
		}
	}
	return id, hashToString(hash), nil
}

// AddReplicaGroup adds a replica group to the hash ring, returning a list of key ranges that need
// to be reassigned to this new group. Specifically, for each new virtual node, the ring must be
// updated and a corresponding reassignment entry must be created (to be returned).
// If the replica group is already in the ring, this method is a no-op, and a nil slice is
// returned.
//
// The reassignment entry for a given virtual node must specify the key range that needs to be
// moved to the new replica group due to the virtual node (and from where). The length of the
// returned list of reassignments must equal the number of virtual nodes per replica group,
// with one entry corresponding to each virtual node (but in any order).

// We construct virtual node ids (which are then hashed, randomizing their spread on the ring)
// by adding a predetermined integer (from 0 to N-1)  to the replica group id.
func (c *ConsistentHash) AddReplicaGroup(id uint64) []Reassignment {

	// TODO(students): [Partitioning] Implement me!
	var res []Reassignment
	for _, value := range c.virtualNodes {
		if value.id == id {
			return res
		}
	}

	nodeGroup := c.virtualNodesForGroup(id)
	c.virtualNodes = append(c.virtualNodes, nodeGroup...)
	slices.SortFunc(c.virtualNodes, virtualNodeLess)
	if len(c.virtualNodes) == c.virtualNodesPerGroup {
		return res
	}
	for i, node := range c.virtualNodes {
		if node.id == id {
			var pre virtualNode
			if i == 0 {
				pre = c.virtualNodes[len(c.virtualNodes)-1]
			} else {
				pre = c.virtualNodes[i-1]
			}
			nextIndex := i
			if nextIndex == len(c.virtualNodes)-1 {
				nextIndex = 0
			} else {
				nextIndex = nextIndex + 1
			}
			for c.virtualNodes[nextIndex].id == id {
				if nextIndex == len(c.virtualNodes)-1 {
					nextIndex = 0
				} else {
					nextIndex = nextIndex + 1
				}
			}
			next := c.virtualNodes[nextIndex]
			keyrange := KeyRange{hashToString(incrementHash(pre.hash)), hashToString(node.hash)}
			res = append(res, Reassignment{next.id, uint64(node.id), keyrange})
		}
	}
	return res
}

// RemoveReplicaGroup removes a replica group from the hash ring, returning a list of key
// ranges that neeed to be reassigned to other replica groups. If the replica group does
// not exist, this method is a no-op, and an empty slice is returned. It is undefined behavior
// to remove the last replica group from the ring, and this will not be tested.
//
// There must be a reassignment entry for every virtual node of the removed group, specifying
// where its keys should be reassigned. The length of the returned list of reassignments must
// equal the number of virtual nodes per replica group (but in any order). The reassignments
// must also account for every key that was previously assigned to the now removed replica group.
func (c *ConsistentHash) RemoveReplicaGroup(id uint64) []Reassignment {

	// TODO(students): [Partitioning] Implement me!
	var res []Reassignment
	nodes := c.virtualNodes
	for i, node := range nodes {
		if node.id == id {
			nextIndex := i
			if nextIndex == len(c.virtualNodes)-1 {
				nextIndex = 0
			} else {
				nextIndex = nextIndex + 1
			}
			for c.virtualNodes[nextIndex].id == id {
				if nextIndex == len(c.virtualNodes)-1 {
					nextIndex = 0
				} else {
					nextIndex = nextIndex + 1
				}
			}
			next := c.virtualNodes[nextIndex]

			var pre virtualNode
			if i == 0 {
				pre = c.virtualNodes[len(c.virtualNodes)-1]
			} else {
				pre = c.virtualNodes[i-1]
			}
			keyrange := KeyRange{hashToString(incrementHash(pre.hash)), hashToString(node.hash)}
			res = append(res, Reassignment{node.id, next.id, keyrange})
		}
	}
	var newNodes []virtualNode
	for _, node := range nodes {
		if node.id != id {
			newNodes = append(newNodes, node)
		}
	}
	c.virtualNodes = newNodes
	return res
}

// ======================================
// DO NOT CHANGE ANY CODE BELOW THIS LINE
// ======================================

// ConsistentHash is a partitioner that implements consistent hashing.
type ConsistentHash struct {
	// virtualNodesPerGroup defines the number of virtual nodes that are created for
	// each replica group.
	virtualNodesPerGroup int

	// virtualNodes defines the hash ring as a sorted list of virtual nodes, starting with the
	// smallest hash value. It must ALWAYS be in ascending sorted order by hash.
	virtualNodes []virtualNode

	// hasher is used to hash all values. Other than pre-defined helpers, this should never be
	// used directly.
	hasher func([]byte) [32]byte
}

// NewConsistentHash creates a new consistent hash partitioner with the default SHA256 hasher.
func NewConsistentHash(virtualNodesPerGroup int) *ConsistentHash {
	return &ConsistentHash{
		virtualNodesPerGroup: virtualNodesPerGroup,
		hasher:               sha256.Sum256,
	}
}

// node returns the virtual node at the specified index.
//
// If the index is out of bounds, it is wrapped using modular arithmetic. For example, an
// index of -1 would map to len(c.virtualNodes)-1.
func (c *ConsistentHash) node(index int) virtualNode {
	clipped := index % len(c.virtualNodes)
	if clipped < 0 {
		clipped += len(c.virtualNodes)
	}
	return c.virtualNodes[clipped]
}

// virtualNodesForGroup returns the virtual nodes for the specified replica group.
// Given the configured parameter, N virtual nodes are created and subsequently returned.
// The virtual nodes are disambiguated by an index that is used when generating their hash.
func (c *ConsistentHash) virtualNodesForGroup(id uint64) []virtualNode {
	var virtualNodes []virtualNode

	for i := 0; i < c.virtualNodesPerGroup; i++ {
		virtualNodeHash := c.virtualNodeHash(id, i)

		virtualNodes = append(virtualNodes, virtualNode{
			id:   id,
			num:  i,
			hash: virtualNodeHash,
		})
	}

	return virtualNodes
}

// virtualNode defines a node in the consistent hash ring. It is a combination
// of the replica group id, the disambiguating virtual number, and the node's hash.
type virtualNode struct {
	id   uint64
	num  int
	hash [32]byte
}

// virtualNodeCmp compares two virtual nodes by their hash, returning -1 if a < b,
// 0 if a == b, and 1 if a > b.
func virtualNodeCmp(a, b virtualNode) int {
	return bytes.Compare(a.hash[:], b.hash[:])
}

// virtualNodeLess compares two virtual nodes by their hash, returning true if and
// only if a < b.
func virtualNodeLess(a, b virtualNode) bool {
	return virtualNodeCmp(a, b) < 0
}

// incrementHash adds 1 to the given hash, wrapping back to 0 if necessary.
func incrementHash(hash [32]byte) [32]byte {
	for i := len(hash) - 1; i >= 0; i-- {
		if hash[i] < math.MaxUint8 {
			hash[i]++
			return hash
		}

		hash[i] = 0
	}
	return hash
}

// hashToString returns the hex string representation of the specified hash. This is useful
// because although we internally represent hashes as byte arrays, we sometimes need to return
// the string hash of a key in our RPC API. It should be used whenever we need to return the
// hash of a key as a string in our API. This includes both specifying reassignemnts and
// creating rewritten keys.
func hashToString(h [32]byte) string {
	return hex.EncodeToString(h[:])
}

// keyHash returns the hash of the specified key.
func (c *ConsistentHash) keyHash(key string) [32]byte {
	hash := c.hasher([]byte(key))

	return hash
}

// virtualNodeHash returns the hash of a virtual node, which is defined by a replica group
// id and number disambiguating different virtual nodes of the same group.
//
// Specifically, the disambiguation number is added to the id before hashing to spread the virtual
// nodes across the ring. Adding, rather than appending, is acceptable since the ids are randomly
// generated and the chance of any conflicts is minimal.
func (c *ConsistentHash) virtualNodeHash(id uint64, virtualNum int) [32]byte {
	virtualID := make([]byte, 8)

	binary.BigEndian.PutUint64(virtualID, id+uint64(virtualNum))

	hash := c.hasher(virtualID)

	return hash
}
