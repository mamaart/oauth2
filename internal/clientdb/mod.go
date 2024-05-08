package clientdb

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type clientDB struct {
	db *gorm.DB
}

func New() *clientDB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(
		&Client{},
		&Scope{},
		&AuthorizationCode{},
	); err != nil {
		panic(err)
	}

	return &clientDB{db: db}
}
