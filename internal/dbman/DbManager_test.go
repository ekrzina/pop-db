package dbman

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/haoli/pop-db/test/mocks"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Fake file structures

type fakeFile struct{}

func newFakeFile() *fakeFile { return &fakeFile{} }

func (f *fakeFile) Read(p []byte) (int, error) { return 0, io.EOF }
func (f *fakeFile) Close() error               { return nil }
func (f *fakeFile) Stat() (fs.FileInfo, error) { return nil, nil }

func setupTestDB(t *testing.T, autoMigrate bool) *DbManager {
	t.Helper()

	dbDir := t.TempDir()
	backupDir := t.TempDir()

	v := viper.New()
	v.Set("database.path", dbDir)
	v.Set("database.name", "test.db")
	v.Set("database.backupPath", backupDir)
	v.Set("database.autoMigrate", autoMigrate)
	v.Set("database.foreignKeys", true)

	logger := zerolog.New(os.Stdout)

	manager, err := NewDbManager(v, logger)
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}

	t.Cleanup(func() {
		if err := manager.DB.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close manager")
		}
	})

	return manager
}

func TestNewDbManager(t *testing.T) {
	t.Run("success with automigrate and foreign keys", func(t *testing.T) {
		tmpDir := t.TempDir()
		viper := viper.New()

		viper.Set("database.path", tmpDir)
		viper.Set("database.name", "test.db")
		viper.Set("database.backupPath", tmpDir)
		viper.Set("database.autoMigrate", true)
		viper.Set("database.foreignKeys", true)
		logger := zerolog.New(os.Stdout)

		manager, err := NewDbManager(viper, logger)
		assert.NoError(t, err)
		assert.NotNil(t, manager)
		assert.NotNil(t, manager.DB)

		// Check that DB file exists
		dbFile := filepath.Join(tmpDir, "test.db")
		_, err = os.Stat(dbFile)
		assert.NoError(t, err)

		// Verify tables were created
		rows, err := manager.DB.Query(`
			SELECT name FROM sqlite_master 
			WHERE type='table' AND name IN ('person','medical_data');
		`)
		assert.NoError(t, err)

		var count int
		for rows.Next() {
			count++
		}
		assert.Equal(t, 2, count)
		assert.NoError(t, manager.DB.Close())
	})
	t.Run("no automigrate does not create tables", func(t *testing.T) {
		tmpDir := t.TempDir()

		v := viper.New()
		v.Set("database.path", tmpDir)
		v.Set("database.name", "test.db")
		v.Set("database.backupPath", tmpDir)
		v.Set("database.autoMigrate", false)
		v.Set("database.foreignKeys", false)

		logger := zerolog.New(os.Stdout)
		manager, err := NewDbManager(v, logger)
		assert.NoError(t, err)

		rows, err := manager.DB.Query(`
		SELECT name FROM sqlite_master 
		WHERE type='table' AND name='person';
	`)
		assert.NoError(t, err)
		assert.False(t, rows.Next())
		assert.NoError(t, manager.DB.Close())
	})
	t.Run("fails if config missing", func(t *testing.T) {
		v := viper.New() // empty config
		logger := zerolog.New(os.Stdout)

		manager, err := NewDbManager(v, logger)
		assert.Error(t, err)
		assert.Nil(t, manager)
	})
}

func TestDbSetupValidate(t *testing.T) {
	t.Run("success on valid config", func(t *testing.T) {
		cfg := DbSetup{
			Path:       "/tmp",
			Name:       "db.sqlite",
			BackupPath: "/backup",
		}
		assert.NoError(t, cfg.validate())
	})

	t.Run("fail on missing path", func(t *testing.T) {
		cfg := DbSetup{
			Name:       "db.sqlite",
			BackupPath: "/backup",
		}
		err := cfg.validate()
		assert.Error(t, err)
		assert.Equal(t, "database.path is required", err.Error())
	})

	t.Run("fail on missing name", func(t *testing.T) {
		cfg := DbSetup{
			Path:       "/tmp",
			BackupPath: "/backup",
		}
		err := cfg.validate()
		assert.Error(t, err)
		assert.Equal(t, "database.name is required", err.Error())
	})

	t.Run("fail on missing backup path", func(t *testing.T) {
		cfg := DbSetup{
			Path: "/tmp",
			Name: "db.sqlite",
		}
		err := cfg.validate()
		assert.Error(t, err)
		assert.Equal(t, "database.backupPath is required", err.Error())
	})
}

func TestMigrate(t *testing.T) {
	manager := setupTestDB(t, false)
	t.Run("successfully creates required tables", func(t *testing.T) {
		err := manager.migrate()
		assert.NoError(t, err)

		rows, err := manager.DB.Query(`
			SELECT name FROM sqlite_master 
			WHERE type='table' AND name IN ('person','medical_data');
		`)
		assert.NoError(t, err)

		var count int
		for rows.Next() {
			count++
		}
		assert.Equal(t, 2, count)
	})
}

func TestExecuteAndQueryWrappers(t *testing.T) {
	manager := setupTestDB(t, true)

	t.Run("success on insert and query person", func(t *testing.T) {
		_, err := manager.DB.Execute(`
			INSERT INTO person (name, surname, nationality, city)
			VALUES (?, ?, ?, ?)`,
			"Arthur", "Dent", "British", "London",
		)
		assert.NoError(t, err)

		rows, err := manager.DB.Query(`
			SELECT name FROM person WHERE surname = ?`,
			"Dent",
		)
		assert.NoError(t, err)
		assert.True(t, rows.Next())

		var name string
		err = rows.Scan(&name)
		assert.NoError(t, err)
		assert.Equal(t, "Arthur", name)
	})

	t.Run("success on queryrow", func(t *testing.T) {
		row := manager.DB.QueryRow(`
			SELECT COUNT(*) FROM person`,
		)

		var count int
		err := row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

func TestBeginTransaction(t *testing.T) {
	manager := setupTestDB(t, true)

	t.Run("success on commit transaction", func(t *testing.T) {
		tx, err := manager.DB.Begin()
		assert.NoError(t, err)

		_, err = tx.Exec(`
			INSERT INTO person (name, surname, nationality, city)
			VALUES (?, ?, ?, ?)`,
			"Ford", "Prefect", "Betelgeusian", "London",
		)
		assert.NoError(t, err)

		assert.NoError(t, tx.Commit())

		row := manager.DB.QueryRow(`SELECT COUNT(*) FROM person`)
		var count int
		assert.NoError(t, row.Scan(&count))
		assert.Equal(t, 1, count)
	})
}

func TestClose(t *testing.T) {
	manager := setupTestDB(t, true)

	err := manager.DB.Close()
	assert.NoError(t, err)

	// Further usage should fail
	_, err = manager.DB.Execute("SELECT 1")
	assert.Error(t, err)
}

func TestGenerateBackupFilename(t *testing.T) {
	manager := setupTestDB(t, false)

	manager.config.Name = "population.db"

	ts := time.Date(2026, 2, 25, 22, 50, 13, 0, time.UTC)

	filename := manager.generateBackupFilename(ts)

	assert.Equal(t,
		"population.db_20260225_225013.bak",
		filename,
	)
}

func TestValidateBackup(t *testing.T) {
	t.Run("successful validation", func(t *testing.T) {
		manager := setupTestDB(t, true)

		tmpFile := filepath.Join(t.TempDir(), "backup.bak")
		err := os.WriteFile(tmpFile, []byte("data"), 0644)
		assert.NoError(t, err)

		err = manager.validatePath(tmpFile)
		assert.NoError(t, err)
	})

	t.Run("fail on not found", func(t *testing.T) {
		manager := setupTestDB(t, true)

		err := manager.validatePath("does_not_exist.bak")
		assert.Error(t, err)
		assert.Equal(t, "path does not exist: does_not_exist.bak", err.Error())
	})
}

func TestWriteBackup(t *testing.T) {
	t.Run("success on writing backup", func(t *testing.T) {
		manager := setupTestDB(t, true)
		mockOS := mocks.NewMockOS(t)

		tmpDir := t.TempDir()
		srcFile, err := os.CreateTemp(tmpDir, "src")
		assert.NoError(t, err)

		dstFile, err := os.CreateTemp(tmpDir, "dst")
		assert.NoError(t, err)

		mockOS.EXPECT().
			Open(mock.Anything).
			Return(srcFile, nil)

		mockOS.EXPECT().
			Create(mock.Anything).
			Return(dstFile, nil)

		mockOS.EXPECT().
			Copy(mock.Anything, mock.Anything).
			Return(int64(123), nil)

		manager.OS = mockOS

		meta, err := manager.WriteBackup()
		assert.NoError(t, err)
		assert.NotNil(t, meta)
		assert.Equal(t, int64(123), meta.SizeBytes)
		assert.NotEmpty(t, meta.Filename)
	})

	t.Run("fail on file open", func(t *testing.T) {
		manager := setupTestDB(t, true)
		mockOS := mocks.NewMockOS(t)

		mockOS.EXPECT().
			Open(mock.Anything).
			Return(nil, fmt.Errorf("open failed"))

		manager.OS = mockOS

		meta, err := manager.WriteBackup()
		assert.Error(t, err)
		assert.Nil(t, meta)
	})

	t.Run("fail on create", func(t *testing.T) {
		manager := setupTestDB(t, true)
		mockOS := mocks.NewMockOS(t)

		tmpFile, _ := os.CreateTemp(t.TempDir(), "src")

		mockOS.EXPECT().
			Open(mock.Anything).
			Return(tmpFile, nil)

		mockOS.EXPECT().
			Create(mock.Anything).
			Return(nil, fmt.Errorf("create failed"))

		manager.OS = mockOS

		meta, err := manager.WriteBackup()
		assert.Error(t, err)
		assert.Nil(t, meta)
	})

	t.Run("fail on copy", func(t *testing.T) {
		manager := setupTestDB(t, true)
		mockOS := mocks.NewMockOS(t)
		tmpDir := t.TempDir()

		srcFile, _ := os.CreateTemp(tmpDir, "src")
		dstFile, _ := os.CreateTemp(tmpDir, "dst")

		mockOS.EXPECT().
			Open(mock.Anything).
			Return(srcFile, nil)

		mockOS.EXPECT().
			Create(mock.Anything).
			Return(dstFile, nil)

		mockOS.EXPECT().
			Copy(mock.Anything, mock.Anything).
			Return(int64(0), fmt.Errorf("copy failed"))

		manager.OS = mockOS

		meta, err := manager.WriteBackup()
		assert.Error(t, err)
		assert.Nil(t, meta)
	})
}

func TestRestoreBackup(t *testing.T) {
	t.Run("success on backup restoration", func(t *testing.T) {
		manager := setupTestDB(t, true)
		mockOS := mocks.NewMockOS(t)
		tmpDir := t.TempDir()

		srcFile, _ := os.CreateTemp(tmpDir, "backup")
		dstFile, _ := os.CreateTemp(tmpDir, "active")

		mockOS.EXPECT().
			Stat(mock.Anything).
			Return(nil, nil)

		mockOS.EXPECT().
			Open(mock.Anything).
			Return(srcFile, nil)

		mockOS.EXPECT().
			Create(mock.Anything).
			Return(dstFile, nil)

		mockOS.EXPECT().
			Copy(mock.Anything, mock.Anything).
			Return(int64(50), nil)

		manager.OS = mockOS

		err := manager.RestoreBackup("file.bak")
		assert.NoError(t, err)
	})

	t.Run("fail on backup missing", func(t *testing.T) {
		manager := setupTestDB(t, true)
		mockOS := mocks.NewMockOS(t)
		mockOS.EXPECT().
			Stat(mock.Anything).
			Return(nil, os.ErrNotExist)

		mockOS.EXPECT().
			IsNotExist(os.ErrNotExist).
			Return(true)

		manager.OS = mockOS

		err := manager.RestoreBackup("missing.bak")
		assert.Error(t, err)
	})

	t.Run("fail on open", func(t *testing.T) {
		manager := setupTestDB(t, true)
		mockOS := mocks.NewMockOS(t)
		backupPath := filepath.Join(manager.config.BackupPath, "backup.bak")

		mockOS.EXPECT().
			Stat(backupPath).
			Return(nil, nil)
		mockOS.EXPECT().
			Open(backupPath).
			Return(nil, errors.New("open failed"))

		manager.OS = mockOS
		err := manager.RestoreBackup("backup.bak")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "open failed")
	})

	t.Run("fail on create", func(t *testing.T) {
		manager := setupTestDB(t, true)
		mockOS := mocks.NewMockOS(t)

		backupPath := filepath.Join(manager.config.BackupPath, "backup.bak")
		destPath := filepath.Join(manager.config.Path, manager.config.Name)

		fakeSrc := newFakeFile()

		mockOS.EXPECT().
			Stat(backupPath).
			Return(nil, nil)
		mockOS.EXPECT().
			Open(backupPath).
			Return(fakeSrc, nil)
		mockOS.EXPECT().
			Create(destPath).
			Return(nil, errors.New("create failed"))

		manager.OS = mockOS

		err := manager.RestoreBackup("backup.bak")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create failed")
	})

	t.Run("fail on copy", func(t *testing.T) {
		manager := setupTestDB(t, true)
		mockOS := mocks.NewMockOS(t)

		backupPath := filepath.Join(manager.config.BackupPath, "backup.bak")
		destPath := filepath.Join(manager.config.Path, manager.config.Name)

		fakeSrc := newFakeFile()
		dstFile, _ := os.CreateTemp(t.TempDir(), "active")

		mockOS.EXPECT().
			Stat(backupPath).
			Return(nil, nil)

		mockOS.EXPECT().
			Open(backupPath).
			Return(fakeSrc, nil)

		mockOS.EXPECT().
			Create(destPath).
			Return(dstFile, nil)

		mockOS.EXPECT().
			Copy(dstFile, fakeSrc). // ← match real call
			Return(int64(0), errors.New("copy failed"))

		manager.OS = mockOS

		err := manager.RestoreBackup("backup.bak")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "copy failed")
	})
}

func TestDeleteBackup(t *testing.T) {
	t.Run("rejects path traversal", func(t *testing.T) {
		manager := setupTestDB(t, true)

		err := manager.DeleteBackup("../evil.bak")
		assert.Error(t, err)
		assert.Equal(t, "invalid backup name", err.Error())
	})

	t.Run("fails when file does not exist", func(t *testing.T) {
		manager := setupTestDB(t, true)

		err := manager.DeleteBackup("missing.bak")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "path does not exist")
	})

	t.Run("successfully deletes existing backup", func(t *testing.T) {
		manager := setupTestDB(t, true)

		tmpDir := manager.config.BackupPath
		filePath := filepath.Join(tmpDir, "test.bak")
		err := os.WriteFile(filePath, []byte("backup"), 0644)
		assert.NoError(t, err)

		err = manager.DeleteBackup("test.bak")
		assert.NoError(t, err)

		_, statErr := os.Stat(filePath)
		assert.True(t, os.IsNotExist(statErr))
	})
}

func TestListBackups(t *testing.T) {
	t.Run("returns empty list when no backups exist", func(t *testing.T) {
		manager := setupTestDB(t, true)

		backups, err := manager.ListBackups()
		assert.NoError(t, err)
		assert.Len(t, backups, 0)
	})

	t.Run("returns backups sorted by CreatedAt descending", func(t *testing.T) {
		manager := setupTestDB(t, true)
		dir := manager.config.BackupPath

		file1 := filepath.Join(dir, "old.bak")
		file2 := filepath.Join(dir, "new.bak")

		assert.NoError(t, os.WriteFile(file1, []byte("old"), 0644))
		assert.NoError(t, os.WriteFile(file2, []byte("new"), 0644))

		oldTime := time.Now().Add(-2 * time.Hour)
		newTime := time.Now()

		assert.NoError(t, os.Chtimes(file1, oldTime, oldTime))
		assert.NoError(t, os.Chtimes(file2, newTime, newTime))

		backups, err := manager.ListBackups()
		assert.NoError(t, err)
		assert.Len(t, backups, 2)

		assert.Equal(t, "new.bak", backups[0].Filename)
		assert.Equal(t, "old.bak", backups[1].Filename)
	})
}
