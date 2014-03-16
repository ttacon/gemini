package gemini

import "errors"

var (
	NoPrimaryKey  = errors.New("table has no primary keys")
	NoDbSpecified = errors.New("cannot perform action without a specified database (currently gemini knows about more/less than one")
	NoDbForStruct = errors.New("no database is specified for the given struct")
)
