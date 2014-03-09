package gemini

type DbDialect int

const (
	MySQL DbDialect = iota
	Postgres
	Sqlite
	MongoDb
)

var TestMode = true
