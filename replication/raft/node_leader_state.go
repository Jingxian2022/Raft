package raft

import (
	"context"
	"encoding/json"
	pb "modist/proto"
	"sync"
)

func (rn *RaftNode) sendHeartbeat() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	followerNodes := rn.node.PeerNodes

	var wg0 sync.WaitGroup
	wg0.Add(len(followerNodes))
	for followerNodeID := range followerNodes {
		go func(id uint64) {
			if id == rn.node.ID {
				wg0.Done()
				return
			}
			defer wg0.Done()
			rn.log.Printf("leader %d is sending heartbeat to %d", rn.node.ID, id)
			conn := rn.node.PeerConns[uint64(id)]
			follower := pb.NewRaftRPCClient(conn)
			_, err := follower.AppendEntries(ctx, &pb.AppendEntriesRequest{
				From:         rn.node.ID,
				To:           uint64(id),
				Term:         rn.GetCurrentTerm(),
				PrevLogIndex: rn.LastLogIndex(),
				PrevLogTerm:  rn.GetLog(rn.LastLogIndex()).GetTerm(),
			})
			if err != nil {
				rn.log.Printf("error sending heartbeat to %d: %v", id, err)
			}
			rn.log.Printf("leader sent heartbeat trying")
		}(followerNodeID)
	}
	wg0.Wait()

	// send heartbeats to all servers periodically
	// heartbeatTicker := time.NewTicker(rn.heartbeatTimeout)
	// defer heartbeatTicker.Stop()

	// var wg sync.WaitGroup
	// for {
	// 	select {
	// 	case <-heartbeatTicker.C:
	// 		for followerNodeID := range followerNodes {
	// 			wg.Add(1)
	// 			go func(id uint64) {
	// 				if id == rn.node.ID {
	// 					wg.Done()
	// 					return
	// 				}
	// 				defer wg.Done()
	// 				rn.log.Printf("leader %d is sending heartbeat to %d", rn.node.ID, id)
	// 				conn := rn.node.PeerConns[uint64(id)]
	// 				follower := pb.NewRaftRPCClient(conn)
	// 				_, err := follower.AppendEntries(ctx, &pb.AppendEntriesRequest{
	// 					From: rn.node.ID,
	// 					To:   uint64(id),
	// 					Term: rn.GetCurrentTerm(),
	// 				})
	// 				if err != nil {
	// 					rn.log.Printf("error sending heartbeat to %d: %v", id, err)
	// 				}
	// 			}(followerNodeID)
	// 		}
	// 		wg.Wait()
	// 	}
	// }
}

func (rn *RaftNode) leaderListen(nextStateC chan stateFunction) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case <-nextStateC:
			return
		case <-rn.appendEntriesC: //???
			rn.log.Printf("leader %d received appendEntriesC", rn.node.ID)
		case <-rn.stopC:
			return
		case proposal := <-rn.proposeC:
			rn.log.Printf("leader %d received proposeC", rn.node.ID)
			rn.HandleProposeC(proposal, nextStateC)
		case voteRequest := <-rn.requestVoteC:
			rn.log.Printf("leader %d received requestVoteC", rn.node.ID)
			if voteRequest.request.Term > rn.GetCurrentTerm() {
				rn.SetCurrentTerm(voteRequest.request.Term)
				nextStateC <- rn.doFollower
				return
			}
		}
	}
}

// leader calls this function
func (rn *RaftNode) HandleProposeC(proposalmsg []byte, nextStateC chan stateFunction) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var proposal pb.PutRequest
	json.Unmarshal(proposalmsg, &proposal)

	append := 1
	success := make(chan bool, 0)

	go func() {
		for {
			if append > len(rn.node.PeerNodes)/2 {
				rn.log.Printf("leader %d received majority of appendEntries", rn.node.ID)
				// commit, then tell others to commit FIXME:should every node commit? or just who received
				leaderCommit := &CommitMsg{
					success:     true,
					err:         nil,
					lastApplied: rn.lastApplied,
					commitIndex: rn.commitIndex,
				}
				tmp, err := json.Marshal(leaderCommit)
				if err != nil {
					rn.log.Printf("error marshalling commit message: %v", err)
				}
				commitMsg := commit(tmp)
				rn.commitC <- &commitMsg
				success <- true
				return
			}
		}
	}()

	for followerNodeID := range rn.node.PeerNodes {
		if followerNodeID == rn.node.ID {
			continue
		}
		go func(id uint64) {
			conn := rn.node.PeerConns[uint64(followerNodeID)]
			follower := pb.NewRaftRPCClient(conn)

			reply, err := follower.AppendEntries(ctx, &pb.AppendEntriesRequest{
				From:         rn.node.ID,
				To:           uint64(id),
				Term:         rn.GetCurrentTerm(),
				PrevLogIndex: rn.LastLogIndex(),
				PrevLogTerm:  rn.GetLog(rn.LastLogIndex()).GetTerm(),
				Entries:      []*pb.LogEntry{{Term: rn.GetCurrentTerm(), Data: proposalmsg}}, // FIXME:
			})
			if err != nil {
				rn.log.Printf("error sending proposal to %d: %v", followerNodeID, err)
			} else {
				if reply.Success {
					append++
				} else {
					if reply.Term > rn.GetCurrentTerm() {
						rn.SetCurrentTerm(reply.Term)
						nextStateC <- rn.doFollower
						return
					}
				}
			}
		}(followerNodeID)
	}

}

// doLeader implements the logic for a Raft node in the leader state.
func (rn *RaftNode) doLeader() stateFunction {
	rn.log.Printf("transitioning to leader state at term %d", rn.GetCurrentTerm())
	rn.state = LeaderState

	// TODO(students): [Raft] Implement me!
	// Hint: perform any initial work, and then consider what a node in the
	// leader state should do when it receives an incoming message on every
	// possible channel.

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	nextStateC := make(chan stateFunction, 1)

	go rn.sendHeartbeat() //TODO: return when leader is not leader anymore!!!
	go rn.leaderListen(nextStateC)

	for {
		select {
		case <-rn.stopC:
			return nil
		case <-nextStateC:
			return <-nextStateC
		}
	}

	return nil
}
