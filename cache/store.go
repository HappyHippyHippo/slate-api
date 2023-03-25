package cache

import (
	"bytes"
	"encoding/gob"
	"time"
)

// IStore is the interface of a cache backend layer.
type IStore interface {
	// Get retrieves an item from the cache. Returns the item or nil, and a
	// bool indicating whether the key was found.
	Get(key string, value interface{}) error

	// Set sets an item to the cache, replacing any existing item.
	Set(key string, value interface{}, expire time.Duration) error

	// Add adds an item to the cache only if an item doesn't already exist
	// for the given key, or if the existing item has expired. Returns
	// an error otherwise.
	Add(key string, value interface{}, expire time.Duration) error

	// Replace sets a new value for the cache key only if it already exists.
	// Returns an error if it does not.
	Replace(key string, data interface{}, expire time.Duration) error

	// Delete removes an item from the cache. Does nothing if the key
	// is not in the cache.
	Delete(key string) error

	// Increment increments a real number, and returns error if the value
	// is not real
	Increment(key string, data uint64) (uint64, error)

	// Decrement decrements a real number, and returns error if the value
	// is not real
	Decrement(key string, data uint64) (uint64, error)

	// Flush sets all items from the cache.
	Flush() error
}

type store struct {
	defaultExpiration time.Duration
}

func (s store) normalizeExpire(
	expire time.Duration,
) time.Duration {
	switch expire {
	case DEFAULT:
		return s.defaultExpiration
	case FOREVER:
		return time.Duration(0)
	}
	return expire
}

func (store) serialize(
	value interface{},
) ([]byte, error) {
	// check if the value can be directly converted into an array of bytes
	if b, ok := value.([]byte); ok {
		return b, nil
	}
	// gob encoding of the data
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	if e := encoder.Encode(value); e != nil {
		return nil, e
	}
	// return the encoding result
	return b.Bytes(), nil
}

func (store) deserialize(
	byt []byte,
	ptr interface{},
) (e error) {
	// check if the given pointer to an array of bytes
	// meaning that can be directly used to store the source byte array
	if b, ok := ptr.(*[]byte); ok {
		*b = byt
		return nil
	}
	// gob decoding of the data
	b := bytes.NewBuffer(byt)
	decoder := gob.NewDecoder(b)
	return decoder.Decode(ptr)
}
