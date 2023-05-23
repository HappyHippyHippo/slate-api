package logmw

import (
	"github.com/happyhippyhippo/slate"
	api "github.com/happyhippyhippo/slate-api"
)

const (
	// ID defines the id to be used as the container
	// registration id of a logging middleware instance factory function.
	ID = api.ID + ".logmw"
)

// Provider defines the slate.rest.log module service provider to be used on
// the application initialization to register the logging middleware service.
type Provider struct{}

var _ slate.IProvider = &Provider{}

// Register will register the log middleware package instances in the
// application container
func (Provider) Register(
	container slate.IContainer,
) error {
	// check container argument reference
	if container == nil {
		return errNilPointer("container")
	}
	_ = container.Service(ID, NewMiddlewareGenerator)
	return nil
}

// Boot will start the migration package
// If the auto migration is defined as true, ether by global variable or
// by environment variable, the migrator will automatically try to migrate
// to the last registered migration
func (p Provider) Boot(
	container slate.IContainer,
) error {
	// check container argument reference
	if container == nil {
		return errNilPointer("container")
	}
	return nil
}
