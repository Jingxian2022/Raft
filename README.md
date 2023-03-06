# Partitioning
## Distribution of Labor
### Yaxin Liu (yliu587)
* Task 1: Lookup

### Jingxian Zhang (jzhan529)
* Task 2: AddReplicaGroup
* Task 3: RemoveReplicaGroup
* Tests

## Test
### TestConsistentHash_Lookup_SimpleIdentity
Testing basic function of Lookup for a ring with two nodes * two groups

### TestConsistentHash_AddReplicaGroup_Basic
Testing basic function of AddReplicaGroup for a ring adding three groups

### TestConsistentHash_AddReplicaGroup_FirstAdd
Testing adding for first time AddReplicaGroup should return nil

### TestConsistentHash_AddReplicaGroup_AlreadyExist
Testing adding for groups that already exist AddReplicaGroup should return nil

### TestConsistentHash_RemoveReplicaGroup
Testing basic function of RemoveReplicaGroup for removing three groups

### TestConsistentHash_RemoveReplicaGroup_NoTargetGroup
Testing adding for groups that not exist RemoveReplicaGroup should return nil