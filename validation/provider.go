package validation

import (
	"github.com/happyhippyhippo/slate"
	api "github.com/happyhippyhippo/slate-api"
)

const (
	// ID defines the id to be used
	// as the container registration id of a validation.
	ID = api.ID + ".validation"

	// TranslatorID defines the id to be used
	// as the container registration id of a translator.
	TranslatorID = ID + ".translator"

	// UniversalTranslatorID defines the id to be used
	// as the container registration id of a universal translator.
	UniversalTranslatorID = TranslatorID + ".universal"

	// ParserID defines the id to be used
	// as the container registration id of an error parser instance.
	ParserID = ID + ".parser"
)

// Provider @todo doc
type Provider struct{}

var _ slate.IProvider = &Provider{}

// Register will register the validation package instances in the
// application container
func (p Provider) Register(
	container slate.IContainer,
) error {
	// check container argument reference
	if container == nil {
		return errNilPointer("container")
	}
	_ = container.Service(TranslatorID, NewTranslator)
	_ = container.Service(UniversalTranslatorID, NewUniversalTranslator)
	_ = container.Service(ParserID, NewParser)
	_ = container.Service(ID, NewValidator)
	return nil
}

// Boot will start the validation package
func (p Provider) Boot(
	container slate.IContainer,
) error {
	// check container argument reference
	if container == nil {
		return errNilPointer("container")
	}
	return nil
}
