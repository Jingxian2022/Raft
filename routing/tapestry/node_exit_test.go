package tapestry

import (
	"context"
	"fmt"
	"testing"
)

func TestLeave(t *testing.T) {
	// Create taps
	tap, _ := MakeTapestries(true, "1", "3", "5", "7") //Make a tapestry with these ids
	fmt.Printf("length of tap %d\n", len(tap))
	tap[1].Leave()
	tap[2].Leave()
	resp, _ := tap[0].FindRoot(
		context.Background(),
		CreateIDMsg("2", 0),
	) //After killing 3 and 5, this should route to 7
	if resp.Next != tap[3].Id.String() {
		t.Errorf("Failed to leave successfully")
	}
}
