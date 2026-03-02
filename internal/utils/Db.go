package utils

import "database/sql"

type DB interface {
	Begin() (*sql.Tx, error)
	Execute(string, ...any) (sql.Result, error)
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
	Close() error
}

// Interface implementation
type RealDB struct {
	DB *sql.DB
}

// Interface compile-time guard
var _ DB = (*RealDB)(nil)

// Begin is a wrapper method for sql.DB's Begin method for transaction support.
// Returns:
//   - *sql.Tx - Tx object for transaction support.
//   - error - Error on transaction beginning setup failure.
func (r *RealDB) Begin() (*sql.Tx, error) {
	return r.DB.Begin()
}

// Execute is a wrapper method of sql.DB's Exec method. A query is executed without returning any rows. Uses background context.
// Parameters:
//   - query: SQL query for database object to execute.
//   - args: Arguments are used for placeholder parameters in the query.
//
// Returns:
//   - sql.Result: The result of the SQL query.
//   - error: An error is returned on query failure.
func (r *RealDB) Execute(query string, args ...any) (sql.Result, error) {
	return r.DB.Exec(query, args...)
}

// Query is a wrapper method of sql.DB's Query method. A query that returns rows is executed. Uses background context.
// Parameters:
//   - query: SQL query for database object to execute.
//   - args: Arguments are used for placeholder parameters in the query.
//
// Returns:
//   - sql.Rows: The result rows of the SQL query.
//   - error: An error is returned on query failure.
func (r *RealDB) Query(query string, args ...any) (*sql.Rows, error) {
	return r.DB.Query(query, args...)
}

// QueryRow is a wrapper method of sql.DB's QueryRow method. A query that returns a row is executed. Uses background context.
// Parameters:
//   - query: SQL query for database object to execute.
//   - args: Arguments are used for placeholder parameters in the query.
//
// Returns:
//   - sql.Rows: The result row of the SQL query.
func (r *RealDB) QueryRow(query string, args ...any) *sql.Row {
	return r.DB.QueryRow(query, args...)
}

// Close is a wrapper for sql.DB's Close method. Closes the database connection.
// Returns:
// - error: Error on closing failure.
func (r *RealDB) Close() error {
	return r.DB.Close()
}
