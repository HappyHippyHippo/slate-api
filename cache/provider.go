package cache

import (
	"github.com/happyhippyhippo/slate"
	api "github.com/happyhippyhippo/slate-api"
)

const (
	// ID defines the id to be used as the container
	// registration id of a cache pool instance, as a base id of all other
	// cache package instances registered in the application container.
	ID = api.ID + ".cache"

	// StoreStrategyTag defines the tag to be assigned to all
	// container Store strategies.
	StoreStrategyTag = ID + ".Store.strategy"

	// StoreFactoryID defines the id to be used as
	//	// the container registration id of a Store factory instance.
	StoreFactoryID = ID + ".Store.factory"
)

// Provider defines the slate.cache module service provider to be used on
// the application initialization to register the caching services.
type Provider struct{}

var _ slate.IProvider = &Provider{}

// Register will register the cache package instances in the
// application container.
func (p Provider) Register(
	container slate.IContainer,
) error {
	// check container argument reference
	if container == nil {
		return errNilPointer("container")
	}
	// add store strategies and factory
	_ = container.Service(StoreFactoryID, NewStoreFactory)
	// add store pool instance
	_ = container.Service(ID, NewStorePool)
	return nil
}

// Boot will start the cache package.
func (p Provider) Boot(
	container slate.IContainer,
) error {
	// check container argument reference
	if container == nil {
		return errNilPointer("container")
	}
	// populate the container Store factory with
	// all registered Store strategies
	storeFactory, e := p.getStoreFactory(container)
	if e != nil {
		return e
	}
	storeStrategies, e := p.getStoreStrategies(container)
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
