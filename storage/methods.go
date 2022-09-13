package storage

type Storage interface {
	Close() error
	AddBackup(backup *Backup) (*Backup, error)
	GetBackup(uuid string) (*Backup, error)
	GetLatestBackup() (*Backup, error)
	DeleteBackup(uuid string) error
	GetBackups() ([]Backup, error)
	AddSettings(body *Settings) (*Settings, error)
	UpdateSettings(uuid string, body *Settings) (*Settings, error)
	DeleteSettings() error
	GetGitToken(uuid string, previewToken bool) (string, error)
	GetSetting(uuid string) (*Settings, error)
}
