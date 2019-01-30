package leveldb

import (
	"sync"

	"github.com/elforg/elfplatform/core/common"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var logger = flogging.MustGetLogger("leveldb")

type dbState int32

const (
	closed dbState = iota
	opened
)

// Config configuration for `DB`
type Config struct {
	DBPath string
}

// DB - a wrapper on an actual store
type DB struct {
	conf    *Config
	db      *leveldb.DB
	dbState dbState
	mutex   sync.Mutex

	readOpts        *opt.ReadOptions
	writeOptsNoSync *opt.WriteOptions
	writeOptsSync   *opt.WriteOptions
}

// CreateDB constructs a `DB`
func CreateDB(conf *Config) *DB {a
	readOpts := &opt.ReadOptions{}
	writeOptsNoSync := &opt.WriteOptions{}
	writeOptsSync := &opt.WriteOptions{}
	writeOptsSync.Sync = true

	return &DB{
		conf:            conf,
		dbState:         closed,
		readOpts:        readOpts,
		writeOptsNoSync: writeOptsNoSync,
		writeOptsSync:   writeOptsSync}
}

// Open opens the underlying db
func (dbInstance *DB) Open() {
	if dbInstance.Opened() {
		return
	}

	dbInstance.mutex.Lock()
	defer dbInstance.mutex.Unlock()
	dbPath := dbInstance.conf.DBPath
	var err error
	var dirEmpty bool
	if dirEmpty, err = createDirAndCheckEmpty(dbPath); err != nil {
		logger.Panicf("Error creating and checking leveldb directory: %s", err)
	}

	dbOpts := &opt.Options{ErrorIfMissing: !dirEmpty}
	if dbInstance.db, err = leveldb.OpenFile(dbPath, dbOpts); err != nil {
		logger.Panicf("Error opening leveldb: %s", err)
	}
	dbInstance.dbState = opened
}

// Opened returns true if dbInstance.dbState has been opened
func (dbInstance *DB) Opened() bool {
	return dbInstance.dbState == opened
}

func createDirAndCheckEmpty(dbPath string) (bool, error) {
	if err := common.CreateDir(dbPath); err != nil {
		logger.Errorf("Error creating leveldb direcotry [%s]", dbPath)
		return false, errors.Wrapf(err, "error creating leveldb direcotry [%s]", dbPath)
	}

	return common.DirEmpty(dbPath)
}

// Close closes the underlying db
func (dbInstance *DB) Close() {
	if !dbInstance.Opened() {
		return
	}

	dbInstance.mutex.Lock()
	defer dbInstance.mutex.Unlock()
	if err := dbInstance.db.Close(); err != nil {
		logger.Errorf("Error closing leveldb: %s", err)
	}
	dbInstance.dbState = closed
}

// Get returns the value for the given key
func (dbInstance *DB) Get(key []byte) ([]byte, error) {
	value, err := dbInstance.db.Get(key, dbInstance.readOpts)
	if err == leveldb.ErrNotFound {
		logger.Warnf("leveldb.ErrNotFound, key: %s", string(key))
		return nil, nil
	}

	if err != nil {
		logger.Errorf("Error retrieving leveldb key [%#v]: %s", key, err)
		return nil, errors.Wrapf(err, "error retrieving leveldb key [%#v]", key)
	}
	return value, nil
}

// Put saves the key/value
func (dbInstance *DB) Put(key []byte, value []byte, sync bool) error {
	writeOpts := dbInstance.makeSyncWriteOptions(sync)
	if err := dbInstance.db.Put(key, value, writeOpts); err != nil {
		logger.Errorf("Error writing leveldb key [%#v]", key)
		return errors.Wrapf(err, "error writing leveldb key [%#v]", key)
	}
	return nil
}

func (dbInstance *DB) makeSyncWriteOptions(sync bool) *opt.WriteOptions {
	if sync {
		return dbInstance.writeOptsSync
	}
	return dbInstance.writeOptsNoSync
}

// Delete deletes the given key
func (dbInstance *DB) Delete(key []byte, sync bool) error {
	writeOpts := dbInstance.makeSyncWriteOptions(sync)
	if err := dbInstance.db.Delete(key, writeOpts); err != nil {
		logger.Errorf("Error deleting leveldb key [%#v]", key)
		return errors.Wrapf(err, "error deleting leveldb key [%#v]", key)
	}
	return nil
}

// GetIterator returns an iterator over key-value store. The iterator should be released after the use.
// The resultset contains all the keys that are present in the db between the startKey (inclusive) and the endKey (exclusive).
// A nil startKey represents the first available key and a nil endKey represent a logical key after the last available key
func (dbInstance *DB) GetIterator(startKey []byte, endKey []byte) iterator.Iterator {
	return dbInstance.db.NewIterator(&util.Range{Start: startKey, Limit: endKey}, dbInstance.readOpts)
}

// WriteBatch writes a batch
func (dbInstance *DB) WriteBatch(batch *leveldb.Batch, sync bool) error {
	writeOpts := dbInstance.makeSyncWriteOptions(sync)
	if err := dbInstance.db.Write(batch, writeOpts); err != nil {
		logger.Errorf("error writing batch to leveldb")
		return errors.Wrap(err, "error writing batch to leveldb")
	}
	return nil
}
