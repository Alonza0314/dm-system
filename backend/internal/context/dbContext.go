package context

import (
	"backend/internal/db"
	"backend/logger"
	"fmt"
	"sync"
)

type dbContextIE struct {
	dbType string
	dbPath string

	*logger.BackendLogger
}

type dbContext struct {
	db db.DbIf

	mainLock    sync.RWMutex
	accountLock sync.RWMutex
	historyLock sync.RWMutex

	*logger.BackendLogger
}

func newDbContext(dbContextIE *dbContextIE) (*dbContext, error) {
	db, err := db.NewDb(dbContextIE.dbType, dbContextIE.dbPath)
	if err != nil {
		return nil, err
	}

	return &dbContext{
		db: db,

		mainLock:    sync.RWMutex{},
		accountLock: sync.RWMutex{},
		historyLock: sync.RWMutex{},

		BackendLogger: dbContextIE.BackendLogger,
	}, nil
}

func (d *dbContext) release() {
	d.DbLog.Infoln("Release dbContext...")

	if err := d.db.Release(); err != nil {
		d.DbLog.Errorf("Failed to release dbContext: %v", err)
	}

	d.DbLog.Infoln("dbContext released")
}

func (d *dbContext) MainLock() *sync.RWMutex {
	return &d.mainLock
}

func (d *dbContext) AccountLock() *sync.RWMutex {
	return &d.accountLock
}

func (d *dbContext) HistoryLock() *sync.RWMutex {
	return &d.historyLock
}

func (d *dbContext) GetHistoryKey(category, device string) string {
	return fmt.Sprintf("%s:%s", category, device)
}

func (d *dbContext) Exist(collection, key string) (bool, error) {
	return d.db.Exist(collection, key)
}

func (d *dbContext) Save(collection, key, value string) error {
	return d.db.Save(collection, key, value)
}

func (d *dbContext) SaveAll(collection string, data map[string]string) error {
	return d.db.SaveAll(collection, data)
}

func (d *dbContext) Load(collection, key string) (string, error) {
	return d.db.Load(collection, key)
}

func (d *dbContext) LoadAll(collection string) (map[string]string, error) {
	return d.db.LoadAll(collection)
}

func (d *dbContext) Remove(collection, key string) error {
	return d.db.Remove(collection, key)
}

func (d *dbContext) RemoveAll(collection string) error {
	return d.db.RemoveAll(collection)
}
