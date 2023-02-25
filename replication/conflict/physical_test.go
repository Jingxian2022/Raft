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

func TestResolveConcurrentEventsConflicts(t *testing.T) {
	r := testCreatePhysicalClockConflictResolver()
	c1 := r.testCreatePhysicalClockGivenTimestamp(20)

	kv1 := KVFromParts[PhysicalClock]("key", "def", c1)
	kv2 := KVFromParts[PhysicalClock]("key", "abc", c1)
	kv, err := r.ResolveConcurrentEvents(kv1, kv2)
	if err != nil {
		t.Errorf("error should be nil")
	}
	if kv.Value != "def" {
		t.Errorf("c1 should happen before c2")
	}
}

func TestResolveConcurrentEventsNoConflicts(t *testing.T) {
	r := testCreatePhysicalClockConflictResolver()

	_, err := r.ResolveConcurrentEvents()
	if err == nil {
		t.Errorf("error should not be nil")
	}
}
