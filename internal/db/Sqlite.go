package db

import (
	"database/sql"
	"path/filepath"
	"pop-db/internal/utils"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// DBSetup structure represents struct to contain database setup data fetched from configuration file
type DBSetup struct {
	Path        string `mapstructure:"path"`
	Name        string `mapstructure:"name"`
	BackupPath  string `mapstructure:"backupPath"`
	AutoMigrate bool   `mapstructure:"autoMigrate"`
	ForeignKeys bool   `mapstructure:"foreignKeys"`
}

// SQLLiteDB structure represents struct to contain database configuration and object
type SQLLiteDB struct {
	cfg    *viper.Viper
	config *DBSetup
	logger *zerolog.Logger
	db     *sql.DB
	OS     utils.OS
}

// migrate creates personal identification and medical table if tables do not already exist.
// Returns:
//   - error: Returned on unsuccessful table creation.
func (s *SQLLiteDB) migrate() error {
	personTable := `
	CREATE TABLE IF NOT EXISTS person (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		surname TEXT NOT NULL,
		occupation TEXT,
		date_of_birth DATETIME,
		nationality TEXT NOT NULL,
		city TEXT NOT NULL,
		notes TEXT,
		picture BLOB,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	medicalTable := `
	CREATE TABLE IF NOT EXISTS medical_data (
		person_id INTEGER PRIMARY KEY,
		height REAL,
		weight REAL,
		blood_type TEXT CHECK(blood_type IN ('A+','A-','B+','B-','AB+','AB-','O+','O-')),
		medical_conditions TEXT,
		FOREIGN KEY(person_id) REFERENCES person(id) ON DELETE CASCADE
	);`
	if _, err := s.db.Exec(personTable); err != nil {
		return err
	}
	if _, err := s.db.Exec(medicalTable); err != nil {
		return err
	}
	return nil
}

// WriteBackup writes backup database to specified configuraion filepath
// Returns:
//   - error: Error on backup fail.
func (s *SQLLiteDB) WriteBackup() error {
	source := filepath.Join(s.config.Path, s.config.Name)
	dest := filepath.Join(
		s.config.BackupPath,
		s.config.Name+"_"+time.Now().Format("20060102_150405")+".bak",
	)
	input, err := s.OS.ReadFile(source)
	if err != nil {
		return err
	}
	return s.OS.WriteFile(dest, input, 0644)
}

// Begin is a wrapper method for sql.DB's Begin method for transaction support.
// Returns:
//   - *sql.Tx - Tx object for transaction support.
//   - error - Error on transaction beginning setup failure.
func (s *SQLLiteDB) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}

// Execute is a wrapper method of sql.DB's Exec method. A query is executed without returning any rows. Uses background context.
// Parameters:
//   - query: SQL query for database object to execute.
//   - args: Arguments are used for placeholder parameters in the query.
//
// Returns:
//   - sql.Result: The result of the SQL query.
//   - error: An error is returned on query failure.
func (s *SQLLiteDB) Execute(query string, args ...any) (sql.Result, error) {
	return s.db.Exec(query, args)
}

// Query is a wrapper method of sql.DB's Query method. A query that returns rows is executed. Uses background context.
// Parameters:
//   - query: SQL query for database object to execute.
//   - args: Arguments are used for placeholder parameters in the query.
//
// Returns:
//   - sql.Rows: The result rows of the SQL query.
//   - error: An error is returned on query failure.
func (s *SQLLiteDB) Query(query string, args ...any) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

// QueryRow is a wrapper method of sql.DB's QueryRow method. A query that returns a row is executed. Uses background context.
// Parameters:
//   - query: SQL query for database object to execute.
//   - args: Arguments are used for placeholder parameters in the query.
//
// Returns:
//   - sql.Rows: The result row of the SQL query.
func (s *SQLLiteDB) QueryRow(query string, args ...any) *sql.Row {
	return s.db.QueryRow(query, args...)
}

// Close is a wrapper for sql.DB's Close method. Closes the database connection.
// Returns:
// - error: Error on closing failure.
func (s *SQLLiteDB) Close() error {
	return s.db.Close()
}

// NewSQLiteDB is a database constructor function.
// Parameters:
//   - v: Viper configuration file for database configuration.
//   - logger: zerolog.Logger object for logging formatting.
//
// Returns:
//   - sql.DB: SQL database structure.
//   - error: Returned on unsuccessful database creation.
func NewSQLiteDB(v *viper.Viper, logger zerolog.Logger) (*SQLLiteDB, error) {
	var cfg DBSetup
	if err := v.UnmarshalKey("database", &cfg); err != nil {
		return nil, err
	}
	fullPath := filepath.Join(cfg.Path, cfg.Name)
	dbase, err := sql.Open("sqlite3", fullPath)
	if err != nil {
		return nil, err
	}
	// Check if foreign keys should be enable enforcing foreign-key constraints
	// SQLLite does not enforce foreign keys by default
	if cfg.ForeignKeys {
		if _, err := dbase.Exec("PRAGMA foreign_keys = ON;"); err != nil {
			return nil, err
		}
	}
	// Create sqllite object
	sqlite := &SQLLiteDB{
		cfg:    v,
		config: &cfg,
		logger: &logger,
		db:     dbase,
		OS:     utils.RealOS{},
	}
	// Check if database should auto migrate
	if cfg.AutoMigrate {
		logger.Info().Msg("Running database migration...")
		if err := sqlite.migrate(); err != nil {
			return nil, err
		}
	}
	return sqlite, nil
}
