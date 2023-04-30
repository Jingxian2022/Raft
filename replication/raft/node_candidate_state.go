package raft

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	pb "modist/proto"
	"sync"
	"time"
)

func (rn *RaftNode) candidateListen(nextStateC chan stateFunction) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// If CommitIndex > lastApplied: increment lastApplied, apply log[lastApplied] to state machine (§5.3)
	// If command received from client: append entry to local log, respond after entry applied to state machine (§5.3)

	for {
		// If RPC request or response contains term T > currentTerm: set currentTerm = T, convert to follower (§5.1)
		select {
		case <-rn.stopC:
			rn.log.Printf("stop message received")
			return
		case <-rn.proposeC:
			candidateCommit := &CommitMsg{
				success:     false,
				err:         errors.New("candidate cannot commit"),
				lastApplied: rn.lastApplied,
				commitIndex: rn.commitIndex,
			}
			tmp, err := json.Marshal(candidateCommit)
			if err != nil {
				rn.log.Printf("error marshalling commit message: %v", err)
				continue
			}
			commitMsg := commit(tmp)
			rn.commitC <- &commitMsg
		// If AppendEntries RPC received from new leader: convert to follower
		case <-rn.appendEntriesC:
			rn.log.Printf("candidate %d received AppendEntries from leader, now convert to follower", rn.node.ID)
			nextStateC <- rn.doFollower
			// TODO: do I need to process the AppendEntries RPC here? or just leave it to the follower state?
			return

		// TODO: check
		case voteReq := <-rn.requestVoteC:
			if voteReq.request.Term > rn.GetCurrentTerm() {
				rn.SetCurrentTerm(voteReq.request.Term)
				nextStateC <- rn.doFollower
				return
			}
		}
	}
}

func (rn *RaftNode) candidateRequestVote(nextStateC chan stateFunction) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Send RequestVote RPCs to all other servers
	fellowNodes := rn.node.PeerNodes
	var wg0 sync.WaitGroup

	// voteReplyC waits for all the vote replies to come back
	voteReplyC := make(chan *pb.RequestVoteReply, len(rn.node.PeerNodes)-1)
	votes := 1
	rej := 0

	go func() {
		for {
			// If votes received from majority of servers: become leader
			if votes > len(rn.node.PeerNodes)/2 {
				nextStateC <- rn.doLeader
				return
			}
			if rej > len(rn.node.PeerNodes)/2 {
				nextStateC <- rn.doFollower
				return
			}
			select {
			case voteR := <-voteReplyC:
				if voteR.Term > rn.GetCurrentTerm() {
					rn.SetCurrentTerm(voteR.Term)
					nextStateC <- rn.doFollower
				}
				if voteR.VoteGranted {
					votes++
				} else {
					rej++
				}
			}
		}
	}()

	for fellowNodeID := range fellowNodes {
		wg0.Add(1)
		go func(id uint64) {
			conn := rn.node.PeerConns[uint64(id)]
			fellow := pb.NewRaftRPCClient(conn)
			voteReply, err := fellow.RequestVote(ctx, &pb.RequestVoteRequest{
				From:         rn.node.ID,
				To:           uint64(id),
				Term:         rn.GetCurrentTerm(),
				LastLogIndex: rn.LastLogIndex(),
				LastLogTerm:  rn.GetLog(rn.LastLogIndex()).GetTerm(),
			})
			if err != nil {
				rn.log.Printf("error sending RequestVote to %d: %v", id, err)
			}
			voteReplyC <- voteReply
		}(fellowNodeID)
	}
	wg0.Wait()

}

// doCandidate implements the logic for a Raft node in the candidate state.
func (rn *RaftNode) doCandidate() stateFunction {
	rn.state = CandidateState
	rn.log.Printf("transitioning to %s state at term %d", rn.state, rn.GetCurrentTerm())

	// TODO(students): [Raft] Implement me!
	// Hint: perform any initial work, and then consider what a node in the
	// candidate state should do when it receives an incoming message on every
	// possible channel.

	nextStateC := make(chan stateFunction, 1)
	go rn.candidateRequestVote(nextStateC)
	go rn.candidateListen(nextStateC)

	// Increment currentTerm
	rn.SetCurrentTerm(rn.GetCurrentTerm() + 1)

	// Vote for self
	rn.setVotedFor(rn.node.ID)

	// Reset election timer, ramdonly choose
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(2) + 1
	timeout := time.After(rn.electionTimeout * time.Duration(randomNum))

	for {
		select {
		// If election timeout elapses: start new election
		case <-timeout:
			return rn.doCandidate

		case nextState := <-nextStateC:
			return nextState
		}
	}
}
