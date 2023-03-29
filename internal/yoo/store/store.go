package store

import (
	"gorm.io/gorm"
	"sync"
)

var (
	once sync.Once
	S    *database
)

type IStore interface {
	Users() UserStore
}

type database struct {
	db *gorm.DB
}

var _IStore = (*database)(nil)

// NewStore returns a new store.
func NewStore(db *gorm.DB) *database {
	once.Do(func() {
		S = &database{db: db}
	})
	return S
}

func (ds *database) Users() UserStore {
	return newUsers(ds.db)
}
