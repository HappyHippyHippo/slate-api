package cache

import (
	"github.com/happyhippyhippo/slate/config"
)

const (
	// UnknownStoreType defines the value to be used to
	// declare an unknown store type.
	UnknownStoreType = "unknown"
)

// IStoreStrategy @todo doc
type IStoreStrategy interface {
	Accept(cfg config.IConfig) bool
	Create(cfg config.IConfig) (IStore, error)
}
