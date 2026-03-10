package dbman

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/haoli/pop-db/internal/utils"

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
	DB     utils.DB
	OS     utils.OS
}

var _ IDbManager = (*DbManager)(nil)

// validate performs non-empty checks on configuration parameters
// Returns:
//   - error: Returned on empty paths in configuration
func (c *DbSetup) validate() error {
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
	path := s.config.Path
	backup := s.config.BackupPath

	// Ensure directories exist
	if err := s.validatePath(path); err != nil {
		if err := s.OS.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	if err := s.validatePath(backup); err != nil {
		if err := s.OS.MkdirAll(backup, os.ModePerm); err != nil {
			return err
		}
	}

	// Create tables
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
	if _, err := s.DB.Execute(personTable); err != nil {
		return err
	}
	if _, err := s.DB.Execute(medicalTable); err != nil {
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

// validatePath checks if database backup is located on designated file
// Parameters:
//   - backupPath: Full file path name of the backup to restore.
//
// Returns:
//   - error: Error if stating file fails
func (s *DbManager) validatePath(path string) error {
	_, err := s.OS.Stat(path)
	if err != nil {
		if s.OS.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}
		return err
	}
	return nil
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
	defer func() {
		if err := sourceFile.Close(); err != nil {
			s.logger.Error().Err(err).Msg("Failed to close source file")
		}
	}()
	// Create destination file and defer closing
	destFile, err := s.OS.Create(destPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := destFile.Close(); err != nil {
			s.logger.Error().Err(err).Msg("Failed to close destination file")
		}
	}()
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

// RestoreBackup overwrites current database with backup database file
// Parameters:
//   - filename: Full file name of the backup to restore to.
//
// Returns:
//   - error: Error returned in case backup does not exist, overwriting the database fails or content is not copied properly.
func (s *DbManager) RestoreBackup(filename string) error {
	backupPath := filepath.Join(s.config.BackupPath, filename)
	if err := s.validatePath(backupPath); err != nil {
		return err
	}
	// Close current connection to database before restoring backup
	realDB, ok := s.DB.(*utils.RealDB)
	if ok && realDB.DB != nil {
		if err := realDB.DB.Close(); err != nil {
			return fmt.Errorf("failed to close active database: %w", err)
		}
	}
	// Open backup file and defer closing
	src, err := s.OS.Open(backupPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := src.Close(); err != nil {
			s.logger.Error().Err(err).Msg("Failed to close source file")
		}
	}()
	// Create new database file by copying backup file to database path
	dstPath := filepath.Join(s.config.Path, s.config.Name)
	dst, err := s.OS.Create(dstPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := dst.Close(); err != nil {
			s.logger.Error().Err(err).Msg("Failed to close destination file")
		}
	}()
	if _, err := s.OS.Copy(dst, src); err != nil {
		return err
	}
	// Reopen database connection after restoring backup
	newDB, err := sql.Open("sqlite3", dstPath)
	if err != nil {
		return fmt.Errorf("failed to reopen database: %w", err)
	}
	if s.config.ForeignKeys {
		if _, err := newDB.Exec("PRAGMA foreign_keys = ON;"); err != nil {
			return fmt.Errorf("failed to re-enable foreign keys: %w", err)
		}
	}
	// Set new database connection
	s.DB = &utils.RealDB{DB: newDB}
	return nil
}

// ListBackups lists all created backups in the config-defined backup directory
// Returns:
//   - []BackupMetadata: A list of backup metadata.
//   - error: Returns error on reading fail.
func (s *DbManager) ListBackups() ([]BackupMetadata, error) {
	entries, err := s.OS.ReadDir(s.config.BackupPath)
	if err != nil {
		return nil, err
	}
	backups := make([]BackupMetadata, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fullPath := filepath.Join(s.config.BackupPath, e.Name())
		info, err := s.OS.Stat(fullPath)
		if err != nil {
			continue
		}
		backups = append(backups, BackupMetadata{
			Filename:  e.Name(),
			Path:      fullPath,
			CreatedAt: info.ModTime(),
			SizeBytes: info.Size(),
		})
	}
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedAt.After(backups[j].CreatedAt)
	})
	return backups, nil
}

// DeleteBackup deletes selected backup from system
// Parameters:
//   - name: Backup file name to delete.
//
// Returns:
//   - error: Error returned on removal fail.
func (s *DbManager) DeleteBackup(name string) error {
	if strings.Contains(name, "..") { // prevents "../" path traversal attacks
		return errors.New("invalid backup name")
	}
	fullPath := filepath.Join(s.config.BackupPath, name)
	if err := s.validatePath(fullPath); err != nil {
		return err
	}
	return s.OS.Remove(fullPath)
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
	if err := cfg.validate(); err != nil {
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
		DB:     &utils.RealDB{DB: dbase},
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
