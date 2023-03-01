# Leaderless Replication
## Distribution of Labor
### Yaxin Liu (yliu587)
* Task 1 Helpers (safelyUpdateKey and getUpToDateKV)
* Task 2 Write Path (HandlePeerWrite, replicateToNode, and ReplicateKey)

### Jingxian Zhang (jzhan529)
* Task 3: Read Path (HandlePeerRead, readFromNode, PerformReadRepair, and GetReplicatedKey)

## Test
### TestHandlePeerRead
Testing basic function of HandlePeerRead after performing a single replicateKey operation

### TestReadRepairAllEventsConcurrent
Testing whether all the nodes get the latest value after performing writing the same key lots of times by performing readRepair correctly 
