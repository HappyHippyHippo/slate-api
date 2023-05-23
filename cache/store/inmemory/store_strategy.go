//go:build inmemory

package inmemory

import (
	"time"

	"github.com/happyhippyhippo/slate-api/cache"
	"github.com/happyhippyhippo/slate/config"
)

const (
	// StoreType defines the value to be used to
	// declare an in-memory store type.
	StoreType = "in-memory"
)

type inMemoryConfig struct {
	DefaultExpiration uint32
}

// StoreStrategy @todo doc
type StoreStrategy struct{}

var _ cache.IStoreStrategy = &StoreStrategy{}

// NewStoreStrategy @todo doc
func NewStoreStrategy() *StoreStrategy {
	return &StoreStrategy{}
}

// Accept @todo doc
func (StoreStrategy) Accept(
	cfg config.IConfig,
) bool {
	// check the config argument reference
	if cfg == nil {
		return false
	}
	// retrieve the data from the configuration
	sc := struct{ Type string }{}
	if _, e := cfg.Populate("", &sc); e != nil {
		return true
	}
	// return acceptance for the read config type
	return sc.Type == StoreType
}

// Create @todo doc
func (StoreStrategy) Create(
	cfg config.IConfig,
) (cache.IStore, error) {
	// check the config argument reference
	if cfg == nil {
		return nil, errNilPointer("config")
	}
	// retrieve the data from the configuration
	sc := inMemoryConfig{
		DefaultExpiration: uint32(cache.DefaultExpiration),
	}
	_, e := cfg.Populate("", &sc)
	if e != nil {
		return nil, e
	}
	// validate configuration
	if sc.DefaultExpiration == 0 {
		return nil, errInvalidStore(cfg, map[string]interface{}{"description": "missing expiration"})
	}
	// return the instantiated in-memory store
	return NewStore(
		time.Duration(sc.DefaultExpiration) * time.Millisecond,
	), nil
}
