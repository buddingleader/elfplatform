package leveldb

// UpdateBatch encloses the details of multiple `updates`
type UpdateBatch struct {
	KVs map[string][]byte
}

// NewUpdateBatch constructs an instance of a Batch
func NewUpdateBatch() *UpdateBatch {
	return &UpdateBatch{make(map[string][]byte)}
}

// Put adds a KV
func (batch *UpdateBatch) Put(key []byte, value []byte) {
	if value == nil {
		panic("Nil value not allowed")
	}
	batch.KVs[string(key)] = value
}

// Delete deletes a Key and associated value
func (batch *UpdateBatch) Delete(key []byte) {
	batch.KVs[string(key)] = nil
}

// Len returns the number of entries in the batch
func (batch *UpdateBatch) Len() int {
	return len(batch.KVs)
}
