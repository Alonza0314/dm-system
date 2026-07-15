package db

import "fmt"

var (
	DefaultDbPath = "/tmp/dm_db.db"
)

var (
	DbTypeList = []string{
		"bbolt",
	}
)

type DbIf interface {
	Release() error

	Save(collection, key, value string) error
	SaveAll(collection string, data map[string]string) error
	Exist(collection, key string) (bool, error)
	Load(collection, key string) (string, error)
	LoadAll(collection string) (map[string]string, error)
	Remove(collection, key string) error
	RemoveAll(collection string) error
}

func NewDb(dbType, dbPath string) (DbIf, error) {
	switch dbType {
	case "bbolt":
		return newBboltDb(dbPath)
	}
	return nil, fmt.Errorf("unsupported db type: %s", dbType)
}