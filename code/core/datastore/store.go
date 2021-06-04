package datastore

import (
	ctx "context"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/0chain/bandwidth_marketplace/code/core/context"
	"github.com/0chain/bandwidth_marketplace/code/core/errors"
	"github.com/0chain/bandwidth_marketplace/code/core/log"
)

const Connection context.CtxKey = "connection"

// Store represent storage type.
type Store struct {
	db *gorm.DB
}

var store Store

// setDB used only for mocking main database.
func setDB(db *gorm.DB) {
	store.db = db
}

// GetStore returns application level storage.
func GetStore() *Store {
	return &store
}

// OpenWithRetries opens database with provided max number of retries.
//
// If an error occurs during execution, the program terminates with code 2 and the last try's
// error will be written in os.Stderr.
//
// OpenWithRetries should be called only once while application starting process.
func (store *Store) OpenWithRetries(cfgDB string, numRetries int) {
	var (
		err               error
		currentNumRetries int
	)
	for currentNumRetries < numRetries {
		err = GetStore().open(cfgDB)
		if err != nil {
			time.Sleep(1 * time.Second)
			currentNumRetries++
			continue
		}
		break
	}

	if err != nil {
		errors.ExitErr("error while open db", err, 2)
	}
}

func (store *Store) open(dsn string) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return DBOpenError
	}

	dbSQL, err := db.DB()
	if err != nil {
		return DBOpenError
	}

	dbSQL.SetMaxIdleConns(100)
	dbSQL.SetMaxOpenConns(200)
	dbSQL.SetConnMaxLifetime(30 * time.Second)
	store.db = db
	return nil
}

// Close closes Store.
func (store *Store) Close() error {
	if store.db != nil {
		if sqldb, _ := store.db.DB(); sqldb != nil {
			if err := sqldb.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

// CreateTransaction begins a transaction and pass pointer
// to it in provided context with Connection key.
func (store *Store) CreateTransaction(cc ctx.Context) ctx.Context {
	db := store.db.Begin()
	// changing type might require further refactor
	return ctx.WithValue(cc, Connection, db) //nolint:staticchec
}

// GetTransaction retrieves transaction from context with Connection key.
func (store *Store) GetTransaction(cc ctx.Context) *gorm.DB {
	conn := cc.Value(Connection)
	if conn != nil {
		return conn.(*gorm.DB)
	}
	log.Logger.Error("No connection in the cc.")
	return nil
}

// GetDB returns current gorm.DB implementation that Store contains.
func (store *Store) GetDB() *gorm.DB {
	return store.db
}
