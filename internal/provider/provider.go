package provider

import (
	"database/sql"
	"log"

	"github.com/spf13/viper"

	"github.com/govinda-attal/cart-commerce/internal/provider/pg"
)

const (
	// PrvDB is used as a key to store db connection within viper config data map.
	PrvDB = "prv.db"
)

// Setup function loads providers for this microservice.
func Setup() {
	db, err := pg.InitStore()
	if err != nil {
		log.Fatal(err)
	}
	viper.SetDefault(PrvDB, db)
}

// DB function returns sql DB connection for this microservice.
func DB() *sql.DB {
	return viper.Get(PrvDB).(*sql.DB)
}

// Cleanup function cleans up active provider resources if any.
// This function is to be called when the microservice is shutting down.
func Cleanup() {
	if db := DB(); db != nil {
		db.Close()
	}
}
