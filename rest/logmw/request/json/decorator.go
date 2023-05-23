package json

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/happyhippyhippo/slate-api/rest/logmw"
	"github.com/happyhippyhippo/slate/log"
)

// NewDecorator will instantiate a new request
// event context reader JSON decorator used to parse the request body as a JSON
// and add the parsed content into the logging data.
func NewDecorator(
	reader logmw.RequestReader,
	model interface{},
) (logmw.RequestReader, error) {
	// check the reader argument reference
	if reader == nil {
		return nil, errNilPointer("reader")
	}
	// return the decorated request reader method
	return func(
		ctx *gin.Context,
	) (log.Context, error) {
		// check the context argument reference
		if ctx == nil {
			return nil, errNilPointer("ctx")
		}
		// read the logging request data from the context
		data, e := reader(ctx)
		if e != nil {
			return nil, e
		}
		// try to unmarshall the request body content if the request
		// is in JSON format, and store it in the data map on the
		// bodyJson field
		contentType := strings.ToLower(ctx.Request.Header.Get("Content-Type"))
		if strings.HasPrefix(contentType, gin.MIMEJSON) {
			if e = json.Unmarshal([]byte(data["body"].(string)), &model); e == nil {
				data["bodyJson"] = model
			}
		}
		// return the request information
		return data, nil
	}, nil
}
