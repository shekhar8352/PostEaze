package database

// Scanner is used to scan the data to the respective types.
type Scanner interface {
	Scan(...interface{}) error
}

// RawEntity is the set of common methods for doing raw queries with an entity.
type RawEntity interface {
	GetQuery(code int) string
	GetQueryValues(code int) []any
	GetMultiQuery(code int) string
	GetMultiQueryValues(code int) []any
	GetNextRaw() RawEntity
	BindRawRow(code int, row Scanner) error
	GetExec(code int) string
	GetExecValues(code int, source string) []any
}

// RawExec is the structure for the entity and the code.
type RawExec struct {
	Entity RawEntity
	Code   int
}
