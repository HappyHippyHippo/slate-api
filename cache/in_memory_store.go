package cache

import (
	"reflect"
	"time"

	"github.com/robfig/go-cache"
)

type inMemoryClient interface {
	Set(key string, value interface{}, expire time.Duration)
	Add(key string, value interface{}, expire time.Duration) error
	Replace(key string, value interface{}, expire time.Duration) error
	Get(key string) (interface{}, bool)
	Increment(key string, delta uint64) (uint64, error)
	Decrement(key string, delta uint64) (uint64, error)
	Delete(key string) bool
	Flush()
}

var _ inMemoryClient = &cache.Cache{}

// InMemoryStore represents the cache with memory persistence
type InMemoryStore struct {
	store
	client inMemoryClient
}

var _ IStore = &InMemoryStore{}

// NewInMemoryStore returns a InMemoryStore
func NewInMemoryStore(
	defaultExpiration time.Duration,
) *InMemoryStore {
	// return the initialized in-memory store struct
	return &InMemoryStore{
		store: store{
			defaultExpiration: defaultExpiration,
		},
		client: cache.New(defaultExpiration, time.Minute),
	}
}

// Get (see IStore interface)
func (c *InMemoryStore) Get(
	key string,
	value interface{},
) error {
	// retrieve the element from the store
	val, found := c.client.Get(key)
	if !found {
		return errMiss(key)
	}
	// try to store the value in the pointer argument
	v := reflect.ValueOf(value)
	if v.Type().Kind() == reflect.Ptr && v.Elem().CanSet() {
		v.Elem().Set(reflect.ValueOf(val))
		return nil
	}
	// signal error while storing the value
	return errNotStored(key)
}

// Set (see IStore interface)
func (c *InMemoryStore) Set(
	key string,
	value interface{},
	expire time.Duration,
) error {
	// store the value in the memory persistence layer
	c.client.Set(key, value, expire)
	return nil
}

// Add (see IStore interface)
func (c *InMemoryStore) Add(
	key string,
	value interface{},
	expire time.Duration,
) error {
	// add the value to the memory, and signal error storing if the
	// key already exists in the memory persistence layer
	err := c.client.Add(key, value, expire)
	if err == cache.ErrKeyExists {
		return errNotStored(key)
	}
	return err
}

// Replace (see IStore interface)
func (c *InMemoryStore) Replace(
	key string,
	value interface{},
	expire time.Duration,
) error {
	// try to replace an existing value in memory
	if err := c.client.Replace(key, value, expire); err != nil {
		return errNotStored(key)
	}
	return nil
}

// Delete (see IStore interface)
func (c *InMemoryStore) Delete(
	key string,
) error {
	// try to remove a value stored in memory marked with the requested key
	if found := c.client.Delete(key); !found {
		return errMiss(key)
	}
	return nil
}

// Increment (see IStore interface)
func (c *InMemoryStore) Increment(
	key string,
	n uint64,
) (uint64, error) {
	// try to increment a value stored in memory or signal a cache miss
	// if not present
	newValue, err := c.client.Increment(key, n)
	if err == cache.ErrCacheMiss {
		return 0, errMiss(key)
	}
	return newValue, err
}

// Decrement (see IStore interface)
func (c *InMemoryStore) Decrement(
	key string,
	n uint64,
) (uint64, error) {
	// try to decrement a value stored in memory or signal a cache miss
	// if not present
	newValue, err := c.client.Decrement(key, n)
	if err == cache.ErrCacheMiss {
		return 0, errMiss(key)
	}
	return newValue, err
}

// Flush (see IStore interface)
func (c *InMemoryStore) Flush() error {
	// flush the cache
	c.client.Flush()
	return nil
}
