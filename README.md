# Conflicts, Consistency, and Clocks
## Test
### VersionVectorClock Test Overview

* TestNewVersionVectorClock: test creating a new VersionVectorClock from VersionVectorConflictResolver
* TestVectorConcurrentEventsDoNotHappenBefore: test whether a VersionVectorClock happens before another
* TestResolveConcurrentEventsConflictsVector: test whether a VersionVectorConflictResolver can merge two concurrent VersionVectorClocks correctly
* TestResolveConcurrentEventsNoConflictsVector: test the rubustness of a VersionVectorConflictResolver when there are no conflicts

