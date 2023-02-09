package tapestry

import (
	"context"
	"fmt"
	"modist/orchestrator/node"
	pb "modist/proto"
	"strconv"
	"strings"
	"sync"
)

type State struct {
	tap *Node
	mu  sync.Mutex

	// The public-facing API that this router must implement
	pb.RouterServer
}

type Args struct {
	Node      *node.Node
	ConnectTo string
}

func Configure(args any) pb.RouterServer {
	a := args.(Args)

	grpcServer := a.Node.GrpcServer
	tap, err := StartTapestry(grpcServer, a.Node.Addr.String(), MakeID(fmt.Sprint(a.Node.ID)), a.ConnectTo)
	if err != nil {
		return nil
	}

	s := &State{
		tap: tap,
	}

	pb.RegisterRouterServer(grpcServer, s)
	go grpcServer.Serve(a.Node.Listener)

	return s
}

/*
Parse an ID from String
*/
func MakeID(stringID string) ID {
	var id ID

	for i := 0; i < DIGITS && i < len(stringID); i++ {
		d, err := strconv.ParseInt(stringID[i:i+1], 16, 0)
		if err != nil {
			return id
		}
		id[i] = Digit(d)
	}
	for i := len(stringID); i < DIGITS; i++ {
		id[i] = Digit(0)
	}

	return id
}

func (s *State) Lookup(ctx context.Context, r *pb.RouteLookupRequest) (*pb.RouteLookupReply, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Note: Separated by comma
	addr, err := s.tap.Get(fmt.Sprint(r.Id))
	if err != nil {
		return nil, err
	}

	addrs := strings.Split(string(addr), ",")

	resp := &pb.RouteLookupReply{
		Addrs: addrs,
	}
	return resp, nil
}

func (s *State) Publish(ctx context.Context, r *pb.RoutePublishRequest) (*pb.RoutePublishReply, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	addr, err := s.tap.Get(fmt.Sprint(r.Id))
	if err != nil {
		addr = []byte{}
	}

	if len(addr) > 0 {
		addr = append(addr, []byte(",")...)
	}

	addr = append(addr, []byte(r.Addr)...)
	err = s.tap.Store(fmt.Sprint(r.Id), []byte(addr))
	if err != nil {
		return nil, err
	}
	return &pb.RoutePublishReply{}, nil
}

func (s *State) Unpublish(ctx context.Context, r *pb.RouteUnpublishRequest) (*pb.RouteUnpublishReply, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	addr, err := s.tap.Get(fmt.Sprint(r.Id))
	if err != nil {
		return nil, err
	}

	addrs := strings.Split(fmt.Sprint(addr), ",")
	for i, a := range addrs {
		if a == r.Addr {
			addrs = append(addrs[:i], addrs[i+1:]...)
			break
		}
	}

	err = s.tap.Store(fmt.Sprint(r.Id), []byte(strings.Join(addrs, ",")))
	if err != nil {
		return nil, err
	}
	return &pb.RouteUnpublishReply{}, nil
}
