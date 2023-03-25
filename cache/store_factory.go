package cache

import (
	"github.com/happyhippyhippo/slate/config"
)

// IStoreFactory defined the interface of a store factory instance.
type IStoreFactory interface {
	Register(strategy IStoreStrategy) error
	Create(cfg config.IConfig) (IStore, error)
}

// StoreFactory is a persistence store generator based on a
// registered list of store generation strategies.
type StoreFactory []IStoreStrategy

var _ IStoreFactory = &StoreFactory{}

// NewStoreFactory @todo doc
func NewStoreFactory() IStoreFactory {
	return &StoreFactory{}
}

// Register will register a new store factory strategy to be used
// on creation requests.
func (f *StoreFactory) Register(
	strategy IStoreStrategy,
) error {
	// check the strategy argument reference
	if strategy == nil {
		return errNilPointer("strategy")
	}
	// add the strategy to the store factory strategy pool
	*f = append(*f, strategy)
	return nil
}

// Create will instantiate and return a new store loaded
// by a configuration instance.
func (f StoreFactory) Create(
	cfg config.IConfig,
) (IStore, error) {
	// check config argument reference
	if cfg == nil {
		return nil, errNilPointer("config")
	}
	// search in the factory strategy pool for one that would accept
	// to generate the requested store with the requested type defined
	// in the given config
	for _, s := range f {
		if s.Accept(cfg) {
			// return the creation of the requested store
			return s.Create(cfg)
		}
	}
	return nil, errInvalidStore(cfg)
}
