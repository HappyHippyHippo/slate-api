package cache

import (
	"time"

	"github.com/happyhippyhippo/slate/config"
)

const (
	// InMemoryStoreType defines the value to be used to
	// declare an in-memory store type.
	InMemoryStoreType = "in-memory"
)

type inMemoryConfig struct {
	DefaultExpiration uint32
}

// InMemoryStoreStrategy @todo doc
type InMemoryStoreStrategy struct{}

var _ IStoreStrategy = &InMemoryStoreStrategy{}

// NewInMemoryStoreStrategy @todo doc
func NewInMemoryStoreStrategy() *InMemoryStoreStrategy {
	return &InMemoryStoreStrategy{}
}

// Accept @todo doc
func (InMemoryStoreStrategy) Accept(
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
	return sc.Type == InMemoryStoreType
}

// Create @todo doc
func (InMemoryStoreStrategy) Create(
	cfg config.IConfig,
) (IStore, error) {
	// check the config argument reference
	if cfg == nil {
		return nil, errNilPointer("config")
	}
	// retrieve the data from the configuration
	sc := inMemoryConfig{
		DefaultExpiration: uint32(DefaultExpiration),
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
	return NewInMemoryStore(
		time.Duration(sc.DefaultExpiration) * time.Millisecond,
	), nil
}
