package response

import (
	"github.com/happyhippyhippo/slate-api/rest/logmw"
)

// ResponseReaderDecorator defines a function used to decorate a response
// reader output.
type ResponseReaderDecorator func(reader logmw.ResponseReader, model interface{}) (logmw.ResponseReader, error)
