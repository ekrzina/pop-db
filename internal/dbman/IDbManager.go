package dbman

type IDbManager interface {
	DeleteBackup(string) error
	ListBackups() ([]BackupMetadata, error)
	RestoreBackup(string) error
	WriteBackup() (*BackupMetadata, error)
}
