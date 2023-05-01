package raft

import (
	"context"
	"encoding/json"
	"math/rand"
	pb "modist/proto"
	"time"
)

func (rn *RaftNode) followerListen(nextState chan stateFunction) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(2) + 1
	timeout := time.After(rn.electionTimeout * time.Duration(randomNum))

	for {
		select {
		case <-timeout:
			rn.log.Printf("follower %d timed out", rn.node.ID)
			nextState <- rn.doCandidate
		case <-rn.stopC:
			rn.log.Printf("stop message received")
			return
		case proposeData := <-rn.proposeC:
			rn.log.Printf("follower %d received proposeC", rn.node.ID)
			// we forward the propose message to the leader
			go rn.handleFollowerProposal(proposeData)

		case voteRequest := <-rn.requestVoteC:
			rn.log.Printf("follower %d received requestVoteC", rn.node.ID)
			if rn.GetVotedFor() == 0 && voteRequest.request.Term >= rn.GetCurrentTerm() {
				rn.setVotedFor(voteRequest.request.From)
				rn.SetCurrentTerm(voteRequest.request.Term)
				// send vote to candidate
				rn.log.Printf("follower %d voted for %d", rn.node.ID, rn.node.ID)
				voteRequest.reply <- pb.RequestVoteReply{
					From:        rn.node.ID,
					To:          voteRequest.request.From,
					Term:        rn.GetCurrentTerm(),
					VoteGranted: true,
				}
			} else {
				// decline the request
				voteRequest.reply <- pb.RequestVoteReply{
					From:        rn.node.ID,
					To:          voteRequest.request.From,
					Term:        rn.GetCurrentTerm(),
					VoteGranted: false,
				}
			}

		case entrymsg := <-rn.appendEntriesC:
			rand.Seed(time.Now().UnixNano())
			randomNum := rand.Intn(2) + 1
			timeout = time.After(rn.electionTimeout * time.Duration(randomNum))

			if entrymsg.request.Term < rn.GetCurrentTerm() {
				// decline the request
				entrymsg.reply <- pb.AppendEntriesReply{
					From:    rn.node.ID,
					To:      entrymsg.request.From,
					Term:    rn.GetCurrentTerm(),
					Success: false,
				}
			} else {
				// accept the request
				// TODO: process the request
				rn.SetCurrentTerm(entrymsg.request.Term)
				rn.leader = entrymsg.request.From
				rn.log.Printf("follower %d received appendEntriesC", rn.node.ID)
				entrymsg.reply <- pb.AppendEntriesReply{
					From:    rn.node.ID,
					To:      rn.node.ID,
					Term:    rn.GetCurrentTerm(),
					Success: true,
				}
				rn.StoreLog(&pb.LogEntry{
					Term:  entrymsg.request.Term,
					Index: entrymsg.request.PrevLogIndex + 1,
				})
			}
		}
	}
}

func (rn *RaftNode) handleFollowerProposal(proposeData []byte) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn := rn.node.PeerConns[rn.leader]
	leader := pb.NewRaftRPCClient(conn)
	// reply := make(chan *pb.ProposalReply)
	// TODO: add for loop to handle timeout
	_, err := leader.Propose(ctx, &pb.ProposalRequest{ // TODO: reply seems contains nothing...
		From: rn.node.ID,
		To:   rn.leader,
		Data: proposeData})
	var promsg pb.PutRequest
	json.Unmarshal(proposeData, promsg)

	if err != nil {
		rn.log.Printf("error sending propose to %d: %v", rn.leader, err)
		tmp, err := json.Marshal(CommitMsg{
			key:     promsg.Key,
			value:   promsg.Value,
			success: false,
		})
		if err != nil {
			rn.log.Printf("error marshalling commit message: %v", err)
		}
		commitMsg := commit(tmp)
		rn.commitC <- &commitMsg
	} else {
		rn.log.Printf("follower %d received proposalReply", rn.node.ID)
		tmp, err := json.Marshal(CommitMsg{
			key:     promsg.Key,
			value:   promsg.Value,
			success: true,
		})
		if err != nil {
			rn.log.Printf("error marshalling commit message: %v", err)
		} else {
			commitMsg := commit(tmp)
			rn.commitC <- &commitMsg
		}
	}
}

// doFollower implements the logic for a Raft node in the follower state.
func (rn *RaftNode) doFollower() stateFunction {
	rn.state = FollowerState
	rn.log.Printf("transitioning to %s state at term %d", rn.state, rn.GetCurrentTerm())

	// TODO(students): [Raft] Implement me!
	// Hint: perform any initial work, and then consider what a node in the
	// follower state should do when it receives an incoming message on every
	// possible channel.
	rn.setVotedFor(0)

	nextState := make(chan stateFunction, 1)
	go rn.followerListen(nextState)

	for {
		select {

		case <-rn.stopC:
			return nil
		case <-nextState:
			return <-nextState
		}
	}
}
