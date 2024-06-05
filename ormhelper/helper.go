package ormhelper

import (
	"database/sql"
	"errors"

	redis "github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	cacheErr "github.com/liyanbing/go-cache/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrDuplicateKey = errors.New("duplicate key")
)

func WrapErr(err error) error {
	if err == nil {
		return nil
	}

	if IsDup(err) {
		return ErrDuplicateKey
	}

	if mongo.IsDuplicateKeyError(err) {
		return ErrDuplicateKey
	}

	switch {
	case errors.Is(err, sql.ErrNoRows),
		errors.Is(err, gorm.ErrRecordNotFound),
		errors.Is(err, mongo.ErrNoDocuments),
		errors.Is(err, mongo.ErrNilDocument),
		errors.Is(err, redis.Nil),
		errors.Is(err, cacheErr.ErrEmptyCache): // nolint
		return ErrNotFound
	}
	return err
}

func IsDup(err error) bool {
	var er *mysql.MySQLError
	if errors.As(err, &er) {
		return er.Number == 1062 // nolint: gomnd
	}
	return false
}
