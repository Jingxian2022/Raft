package conflict

import (
	"testing"
)

func TestNewVersionVectorClock(t *testing.T) {
	r := NewVersionVectorConflictResolver()
	r.vector[0] = 1
	r.vector[1] = 2

	newClock := r.NewClock()
	if newClock.vector[0] != 1 || newClock.vector[1] != 2 {
		t.Errorf("VersionVectorClock's initialization is incorrect")
	}
}

func TestVectorConcurrentEventsDoNotHappenBefore(t *testing.T) {
	v1 := NewVersionVectorClock()
	v2 := NewVersionVectorClock()

	v1.vector[0] = 0
	v1.vector[3] = 2
	v2.vector[0] = 1

	if v2.HappensBefore(v1) {
		t.Errorf("v2 does not happen before v1 due to v1[3]")
	}
	if v1.HappensBefore(v2) {
		t.Errorf("v1 does not happen before v2 due to v2[0] > v1[0]")
	}
}

func TestResolveConcurrentEventsConflictsVector(t *testing.T) {
	r := NewVersionVectorConflictResolver()

	v1 := NewVersionVectorClock()
	v2 := NewVersionVectorClock()

	v1.vector[0] = 0
	v1.vector[3] = 2
	v2.vector[0] = 1

	kv1 := KVFromParts[VersionVectorClock]("key", "value1", v1)
	kv2 := KVFromParts[VersionVectorClock]("key", "value2", v2)

	kv, err := r.ResolveConcurrentEvents(kv1, kv2)
	if err != nil {
		t.Errorf("error should be nil")
	}
	if kv.Value != "value2" {
		t.Errorf("Resolver should choose the key-value with the highest lexicographic value")
	}
	if kv.Clock.vector[0] != 1 || kv.Clock.vector[3] != 2 {
		t.Errorf("Resolver should return a higher clock")
	}
}

func TestResolveConcurrentEventsNoConflictsVector(t *testing.T) {
	r := NewVersionVectorConflictResolver()

	_, err := r.ResolveConcurrentEvents()
	if err == nil {
		t.Errorf("error should not be nil")
	}
}
