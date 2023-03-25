package cache

import (
	"fmt"
	"io"

	"github.com/happyhippyhippo/slate/config"
)

// IStorePool defines the interface of a store pool instance.
type IStorePool interface {
	Get(name string) (IStore, error)
}

// storePool is a database store pool and generator.
type storePool struct {
	cfg          config.IManager
	storeFactory IStoreFactory
	instances    map[string]IStore
}

var _ IStorePool = &storePool{}

// NewStorePool will instantiate a new relational
// database store pool instance.
func NewStorePool(
	cfg config.IManager,
	factory IStoreFactory,
) (IStorePool, error) {
	// check config argument reference
	if cfg == nil {
		return nil, errNilPointer("config")
	}
	// check storeFactory argument reference
	if factory == nil {
		return nil, errNilPointer("factory")
	}
	// instantiate the store pool instance
	pool := &storePool{
		cfg:          cfg,
		storeFactory: factory,
		instances:    map[string]IStore{},
	}
	// check if is to observe store configuration changes
	if ObserveConfig {
		// add an observer to the stores config
		_ = cfg.AddObserver(StoresConfigPath, func(_ interface{}, _ interface{}) {
			// close all the currently opened stores
			for _, store := range pool.instances {
				if c, ok := store.(io.Closer); ok {
					_ = c.Close()
				}
			}
			// clear the storing pool
			pool.instances = map[string]IStore{}
		})
	}
	return pool, nil
}

// Get execute the process of the store creation based on the
// base configuration defined by the given name of the store,
// and apply the extra store cfg also given as arguments.
func (f *storePool) Get(
	name string,
) (IStore, error) {
	// check if the store as already been created and return it
	if store, ok := f.instances[name]; ok {
		return store, nil
	}
	// generate the configuration path of the requested store
	path := fmt.Sprintf("%s.%s", StoresConfigPath, name)
	// check if there is a configuration for the requested store
	if !f.cfg.Has(path) {
		return nil, errConfigNotFound(path)
	}
	// obtain the store configuration
	cfg, e := f.cfg.Config(path)
	if e != nil {
		return nil, e
	}
	// create the store
	store, e := f.storeFactory.Create(cfg)
	if e != nil {
		return nil, e
	}
	// store the store instance
	f.instances[name] = store
	return store, nil
}
