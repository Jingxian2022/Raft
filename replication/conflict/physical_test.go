package conflict

import (
	"testing"
	"time"
)

func testCreatePhysicalClockConflictResolver() *PhysicalClockConflictResolver {
	return &PhysicalClockConflictResolver{}
}

func (c *PhysicalClockConflictResolver) testCreatePhysicalClock() PhysicalClock {
	return PhysicalClock{timestamp: uint64(time.Now().UnixNano())}
}

func (c *PhysicalClockConflictResolver) testCreatePhysicalClockGivenTimestamp(timestamp uint64) PhysicalClock {
	return PhysicalClock{timestamp: timestamp}
}

func TestPhysicalConcurrentEventsHappenBefore(t *testing.T) {
	r := testCreatePhysicalClockConflictResolver()
	c1 := r.testCreatePhysicalClockGivenTimestamp(10)
	c2 := r.testCreatePhysicalClockGivenTimestamp(20)
	c3 := r.testCreatePhysicalClockGivenTimestamp(30)

	if !c1.HappensBefore(c2) {
		t.Errorf("c1 should happen before c2")
	}
	if !c1.HappensBefore(c3) {
		t.Errorf("c1 should happen before c3")
	}
	if !c2.HappensBefore(c3) {
		t.Errorf("c2 should happen before c3")
	}
}
