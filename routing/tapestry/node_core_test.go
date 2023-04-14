package tapestry

import (
	"fmt"
	"testing"
)

// func TestRegister(t *testing.T) {
// 	tap, _ := MakeTapestries(true, "1", "3", "5", "7") //Make a tapestry with these ids
// 	fmt.Printf("length of tap %d\n", len(tap))
// 	// register a node
// 	tap[0].Register(context.Background(), &pb.Registration{FromNode: "5", Key: "a"})
// 	nodes := tap[0].LocationsByKey.Get("a")
// 	fmt.Printf("nodes %v", nodes)
// 	if nodes[0] != tap[0].Id {
// 		t.Errorf("Failed to register successfully")
// 	}
// }

func TestPublish(t *testing.T) {
	tap, _ := MakeTapestries(true, "1", "3", "5", "7") //Make a tapestry with these ids
	fmt.Printf("length of tap %d\n", len(tap))
	// publish a message
	tap[0].Publish("a")
	nodes := tap[0].LocationsByKey.Get("a")
	if nodes[0] != MakeID(1) {
		t.Errorf("Failed to publish successfully")
	}
}

func TestLookup(t *testing.T) {
	tap, _ := MakeTapestries(true, "1000", "1110", "5", "7") //Make a tapestry with these ids
	fmt.Printf("length of tap %d\n", len(tap))
	// publish a message
	tap[0].Publish("110")
	nodes, err := tap[3].Lookup("110")
	if nodes[0] != tap[0].Id && err == nil {
		t.Errorf("Failed to lookup successfully")
	}
}
