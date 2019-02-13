package leveldb

import (
	"bytes"
	"sync"
)

var dbNameKeySeparator = []byte{0x00}
var lastKeyIndicator = byte(0x01)

// Provider enables to use a single leveldb as multiple logical leveldbs
type Provider struct {
	db        *DB
	dbHandles map[string]*DBHandle
	mutex     sync.Mutex
}

// NewProvider constructs a Provider
func NewProvider(conf *Config) *Provider {
	db := CreateDB(conf)
	db.Open()
	return &Provider{db, make(map[string]*DBHandle), sync.Mutex{}}
}

// GetDBHandle returns a handle to a named db
func (p *Provider) GetDBHandle(dbName string) *DBHandle {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	dbHandle := p.dbHandles[dbName]
	if dbHandle == nil {
		dbHandle = &DBHandle{dbName, p.db}
		p.dbHandles[dbName] = dbHandle
	}
	return dbHandle
}

// Close closes the underlying leveldb
func (p *Provider) Close() {
	p.db.Close()
}

func constructLevelKey(dbName string, key []byte) []byte {
	return append(append([]byte(dbName), dbNameKeySeparator...), key...)
}

func retrieveAppKey(levelKey []byte) []byte {
	return bytes.SplitN(levelKey, dbNameKeySeparator, 2)[1]
}
