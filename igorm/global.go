package igorm

import (
	"database/sql"
	"gorm.io/gorm"
)

var (
	globalDB   *sql.DB
	globalGORM *gorm.DB
)

func Close() error {
	if globalDB != nil {
		return globalDB.Close()
	}
	return nil
}

func Session() *gorm.DB {
	return globalGORM
}

func Begin(opts ...*sql.TxOptions) *gorm.DB {
	return globalGORM.Begin(opts...)
}
