package db

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"pop-db/internal/utils"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// DBSetup structure represents struct to contain database setup data fetched from configuration file
type DbSetup struct {
	Path        string `mapstructure:"path"`
	Name        string `mapstructure:"name"`
	BackupPath  string `mapstructure:"backupPath"`
	AutoMigrate bool   `mapstructure:"autoMigrate"`
	ForeignKeys bool   `mapstructure:"foreignKeys"`
}

// BackupMetadata represents additional information on created backup
type BackupMetadata struct {
	Filename  string
	Path      string
	CreatedAt time.Time
	SizeBytes int64
}

// DbManager structure represents struct to contain database configuration and object
type DbManager struct {
	cfg    *viper.Viper
	config *DbSetup
	logger *zerolog.Logger
	db     *sql.DB
	OS     utils.OS
}

// Validate performs non-empty checks on configuration parameters
// Returns:
//   - error: Returned on empty paths in configuration
func (c *DbSetup) Validate() error {
	if c.Path == "" {
		return fmt.Errorf("database.path is required")
	}
	if c.Name == "" {
		return fmt.Errorf("database.name is required")
	}
	if c.BackupPath == "" {
		return fmt.Errorf("database.backupPath is required")
	}
	return nil
}

// migrate creates personal identification and medical table if tables do not already exist.
// Returns:
//   - error: Returned on unsuccessful table creation.
func (s *DbManager) migrate() error {
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

// generateBackupFilename generates new backup file name from current time and configuration parameters
// Parameters:
//   - t: Time of database backup writing request.
func (s *DbManager) generateBackupFilename(t time.Time) string {
	return s.config.Name + "_" + t.Format("20060102_150405") + ".bak"
}

// WriteBackup writes backup database to specified configuraion filepath
// Returns:
//   - *BackupMetadata: metadata of created backup file
//   - error: Error on backup fail.
func (s *DbManager) WriteBackup() (*BackupMetadata, error) {
	now := time.Now().UTC()

	sourcePath := filepath.Join(s.config.Path, s.config.Name)
	filename := s.generateBackupFilename(now)
	destPath := filepath.Join(s.config.BackupPath, filename)
	// Open source file database and defer closing
	sourceFile, err := s.OS.Open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer sourceFile.Close()
	// Create destination file and defer closing
	destFile, err := s.OS.Create(destPath)
	if err != nil {
		return nil, err
	}
	defer destFile.Close()
	// Stream copy of source file to destination file
	size, err := s.OS.Copy(destFile, sourceFile)
	if err != nil {
		return nil, err
	}
	metadata := &BackupMetadata{
		Filename:  filename,
		Path:      destPath,
		CreatedAt: now,
		SizeBytes: size,
	}
	return metadata, nil
}

// ValidateBackup checks if database backup is located on designated file
// Parameters:
//   - backupPath: Full file path name of the backup to restore.
//
// Returns:
//   - error: Error if stating file fails
func (s *DbManager) validateBackup(backupPath string) error {
	if _, err := s.OS.Stat(backupPath); err != nil {
		if s.OS.IsNotExist(err) {
			return fmt.Errorf("backup not found")
		}
		return err
	}
	return nil
}

// RestoreBackup overwrites current database with backup database file
// Parameters:
//   - filename: Full file name of the backup to restore to.
//
// Returns:
//   - error: Error returned in case backup does not exist, overwriting the database fails or content is not copied properly.
func (s *DbManager) RestoreBackup(filename string) error {
	backupPath := filepath.Join(s.config.BackupPath, filename)
	if err := s.validateBackup(backupPath); err != nil {
		return err
	}
	// Open backup file and defer closing
	src, err := s.OS.Open(backupPath)
	if err != nil {
		return err
	}
	defer src.Close()
	// Create/overwrite active database file
	dst, err := s.OS.Create(filepath.Join(s.config.Path + s.config.Name))
	if err != nil {
		return err
	}
	defer dst.Close()
	// Copy contents
	if _, err := s.OS.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

// Begin is a wrapper method for sql.DB's Begin method for transaction support.
// Returns:
//   - *sql.Tx - Tx object for transaction support.
//   - error - Error on transaction beginning setup failure.
func (s *DbManager) Begin() (*sql.Tx, error) {
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
func (s *DbManager) Execute(query string, args ...any) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

// Query is a wrapper method of sql.DB's Query method. A query that returns rows is executed. Uses background context.
// Parameters:
//   - query: SQL query for database object to execute.
//   - args: Arguments are used for placeholder parameters in the query.
//
// Returns:
//   - sql.Rows: The result rows of the SQL query.
//   - error: An error is returned on query failure.
func (s *DbManager) Query(query string, args ...any) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

// QueryRow is a wrapper method of sql.DB's QueryRow method. A query that returns a row is executed. Uses background context.
// Parameters:
//   - query: SQL query for database object to execute.
//   - args: Arguments are used for placeholder parameters in the query.
//
// Returns:
//   - sql.Rows: The result row of the SQL query.
func (s *DbManager) QueryRow(query string, args ...any) *sql.Row {
	return s.db.QueryRow(query, args...)
}

// Close is a wrapper for sql.DB's Close method. Closes the database connection.
// Returns:
// - error: Error on closing failure.
func (s *DbManager) Close() error {
	return s.db.Close()
}

// NewDbManager is a database constructor function.
// Parameters:
//   - v: Viper configuration file for database configuration.
//   - logger: zerolog.Logger object for logging formatting.
//
// Returns:
//   - sql.DB: SQL database structure.
//   - error: Returned on unsuccessful database creation.
func NewDbManager(v *viper.Viper, logger zerolog.Logger) (*DbManager, error) {
	var cfg DbSetup
	if err := v.UnmarshalKey("database", &cfg); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
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
	sqlite := &DbManager{
		cfg:    v,
		config: &cfg,
		logger: &logger,
		db:     dbase,
		OS:     &utils.RealOS{},
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
