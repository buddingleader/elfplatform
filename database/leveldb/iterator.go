package leveldb

import "github.com/syndtr/goleveldb/leveldb/iterator"

// Iterator extends actual leveldb iterator
type Iterator struct {
	iterator.Iterator
}

// Key wraps actual leveldb iterator method
func (itr *Iterator) Key() []byte {
	return retrieveAppKey(itr.Iterator.Key())
}
