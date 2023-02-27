# Conflicts, Consistency, and Clocks
## Distribution of Labor
### Yaxin Liu (yliu587)
* VersionVectorClock implement
* VersionVectorClock tests

### Jingxian Zhang (jzhan529)
* PhysicalClock implement
* PhysicalClock tests

## Test
### VersionVectorClock Test Overview

* TestNewVersionVectorClock: test creating a new VersionVectorClock from VersionVectorConflictResolver
* TestVectorConcurrentEventsDoNotHappenBefore: test whether a VersionVectorClock happens before another
* TestResolveConcurrentEventsConflictsVector: test whether a VersionVectorConflictResolver can merge two concurrent VersionVectorClocks correctly
* TestResolveConcurrentEventsNoConflictsVector: test the rubustness of a VersionVectorConflictResolver when there are no conflicts

### PhysicalClock Test Overview

* TestPhysicalConcurrentEventsHappenBefore: test whether returns the correct happen before judgement with different PhysicalClocks
* TestResolveConcurrentEventsConflicts: test solution of eventual consistency when there is conflict of KVs with the same PhysicalClock
* TestResolveConcurrentEventsNoConflicts: test the rubustness of a PhysicalClockConflictResolver when there are no conflicts
