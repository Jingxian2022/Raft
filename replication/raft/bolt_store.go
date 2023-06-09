package raft

import (
	pb "modist/proto"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

var (
	dbLogs  = []byte("logs")
	dbState = []byte("state")
	dbSnap  = []byte("snapshots")
)

const dbFileMode = 0666

// BoltStore implements the StableStore interface and serves as a storage option for raft
type BoltStore struct {
	db   *bolt.DB
	path string
}

// NewBoltStore either creates a new store or reopens an existing store.
// Also populates the latest logIndex into memory.
func NewBoltStore(path string) *BoltStore {
	db, err := bolt.Open(path, dbFileMode, &bolt.Options{
		Timeout: 2 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		getBucket(tx, dbLogs)
		getBucket(tx, dbState)
		getBucket(tx, dbSnap)
		return nil
	})
	return &BoltStore{db: db, path: path}
}

func getBucket(tx *bolt.Tx, bucket []byte) *bolt.Bucket {
	var err error
	b := tx.Bucket(bucket)
	if b == nil {
		b, err = tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			panic(err)
		}
	}
	return b
}

// SetBytes sets a key-value pair into Bolt
func (store *BoltStore) SetBytes(key, value []byte) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		b := getBucket(tx, dbState)
		return b.Put(key, value)
	})
}

// GetBytes retrives a value for the specified key from Bolt
func (store *BoltStore) GetBytes(key []byte) []byte {
	tx, err := store.db.Begin(false)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()
	b := getBucket(tx, dbState)
	return b.Get(key)
}

// SetUint64 sets a key-value pair into Bolt where the value is a uint64
func (store *BoltStore) SetUint64(key []byte, term uint64) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		b := getBucket(tx, dbState)
		return b.Put(key, uint64ToBytes(term))
	})
}

// GetUint64 retrieves a uint64 from Bolt
func (store *BoltStore) GetUint64(key []byte) uint64 {
	tx, err := store.db.Begin(false)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()
	b := getBucket(tx, dbState)
	value := b.Get(key)
	if value == nil {
		return 0
	}
	return bytesToUint64(value)
}

// StoreLog grabs the next log index and stores a LogEntry into Bolt
func (store *BoltStore) StoreLog(log *pb.LogEntry) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		b := getBucket(tx, dbLogs)
		buf, err := encodeMsgPack(log)
		if err != nil {
			return err
		}
		return b.Put(uint64ToBytes(log.Index), buf.Bytes())
	})
}

// GetLog retrieves a LogEntry at a specific log index from Bolt
func (store *BoltStore) GetLog(index uint64) *pb.LogEntry {
	tx, err := store.db.Begin(false)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()
	b := getBucket(tx, dbLogs)
	value := b.Get(uint64ToBytes(index))
	if value == nil {
		return nil
	}
	var log pb.LogEntry
	err = decodeMsgPack(value, &log)
	if err != nil {
		panic(err)
	}
	return &log
}

// LastLogIndex gets the last index inserted into Bolt
func (store *BoltStore) LastLogIndex() uint64 {
	tx, err := store.db.Begin(false)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()
	b := getBucket(tx, dbLogs)
	cursor := b.Cursor()
	key, _ := cursor.Last()
	if key == nil {
		return 0
	}
	return bytesToUint64(key)
}

// TruncateLog deletes all logs starting from index
func (store *BoltStore) TruncateLog(index uint64) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		b := getBucket(tx, dbLogs)
		cursor := b.Cursor()
		for key, _ := cursor.Seek(uint64ToBytes(index)); key != nil; key, _ = cursor.Next() {
			err := cursor.Delete()
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// AllLogs returns all logs in ascending order. Used for testing purposes.
func (store *BoltStore) AllLogs() []*pb.LogEntry {
	tx, err := store.db.Begin(false)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()
	b := getBucket(tx, dbLogs)
	cursor := b.Cursor()
	var result []*pb.LogEntry
	for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
		var log pb.LogEntry
		if err := decodeMsgPack(value, &log); err != nil {
			panic(err)
		}
		result = append(result, &log)
	}
	return result
}

// Close releases the lock on the db file and closes Bolt
func (store *BoltStore) Close() {
	err := store.db.Close()
	if err != nil && err != bolt.ErrDatabaseNotOpen {
		panic(err)
	}
}

// Remove deletes the db file
func (store *BoltStore) Remove() {
	err := os.RemoveAll(store.path)
	if err != nil {
		panic(err)
	}
}

// Path returns the db path
func (store *BoltStore) Path() string {
	return store.path
}
