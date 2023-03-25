package validation

import (
	api "github.com/happyhippyhippo/slate-api"
	"github.com/happyhippyhippo/slate/env"
)

const (
	// EnvID defines the slate.api.validation package base environment variable name.
	EnvID = api.EnvID + "_VALIDATION"
)

var (
	// Locale defines the default locale string to be used when
	// instantiating the translator.
	Locale = env.String(EnvID+"_LOCALE", "en")
)
