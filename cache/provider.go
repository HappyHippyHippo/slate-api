package cache

import (
	"github.com/happyhippyhippo/slate"
)

const (
	// ID defines the id to be used as the container
	// registration id of a cache pool instance, as a base id of all other
	// cache package instances registered in the application container.
	ID = slate.ID + ".cache"

	// StoreStrategyTag defines the tag to be assigned to all
	// container store strategies.
	StoreStrategyTag = ID + ".store.strategy"

	// InMemoryStrategyID defines the id to be used as
	// the container registration id of an in-memory store factory
	// strategy instance.
	InMemoryStrategyID = ID + ".store.strategy.in_memory"

	// MemcachedStrategyID defines the id to be used as
	// the container registration id of a memcached service store factory
	// strategy instance.
	MemcachedStrategyID = ID + ".store.strategy.memcached"

	// BinaryMemcachedStrategyID defines the id to be used as
	// the container registration id of a binary connection memcached
	// service store factory strategy instance.
	BinaryMemcachedStrategyID = ID + ".store.strategy.binary_memcached"

	// RedisStrategyID defines the id to be used as
	// the container registration id of a redis
	// service store factory strategy instance.
	RedisStrategyID = ID + ".store.strategy.redis"

	// StoreFactoryID defines the id to be used as
	//	// the container registration id of a store factory instance.
	StoreFactoryID = ID + ".store.factory"
)

// Provider defines the slate.cache module service provider to be used on
// the application initialization to register the caching services.
type Provider struct{}

var _ slate.IProvider = &Provider{}

// Register will register the cache package instances in the
// application container.
func (p Provider) Register(
	container ...slate.IContainer,
) error {
	// check container argument reference
	if len(container) == 0 || container[0] == nil {
		return errNilPointer("container")
	}
	// add store strategies and factory
	_ = container[0].Service(InMemoryStrategyID, NewInMemoryStoreStrategy, StoreStrategyTag)
	_ = container[0].Service(StoreFactoryID, NewStoreFactory)
	// add store pool instance
	_ = container[0].Service(ID, NewStorePool)
	return nil
}

// Boot will start the cache package.
func (p Provider) Boot(
	container ...slate.IContainer,
) error {
	// check container argument reference
	if len(container) == 0 || container[0] == nil {
		return errNilPointer("container")
	}
	// populate the container store factory with
	// all registered store strategies
	storeFactory, e := p.getStoreFactory(container[0])
	if e != nil {
		return e
	}
	storeStrategies, e := p.getStoreStrategies(container[0])
	if e != nil {
		return e
	}
	for _, strategy := range storeStrategies {
		_ = storeFactory.Register(strategy)
	}
	return nil
}

func (Provider) getStoreFactory(
	container slate.IContainer,
) (IStoreFactory, error) {
	// retrieve the factory entry
	entry, e := container.Get(StoreFactoryID)
	if e != nil {
		return nil, e
	}
	// validate the retrieved entry type
	instance, ok := entry.(IStoreFactory)
	if !ok {
		return nil, errConversion(entry, "cache.IStoreFactory")
	}
	return instance, nil
}

func (Provider) getStoreStrategies(
	container slate.IContainer,
) ([]IStoreStrategy, error) {
	// retrieve the strategies entries
	entries, e := container.Tag(StoreStrategyTag)
	if e != nil {
		return nil, e
	}
	// type check the retrieved strategies
	var strategies []IStoreStrategy
	for _, entry := range entries {
		if instance, ok := entry.(IStoreStrategy); ok {
			strategies = append(strategies, instance)
		}
	}
	return strategies, nil
}
