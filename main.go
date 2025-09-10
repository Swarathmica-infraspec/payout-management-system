package main

import (
	"database/sql"
	"log"
	"os"
	payee "payoutmanagementsystem/payee"
)

var store *payee.PayeePostgresDB

func initStore() *payee.PayeePostgresDB {
	if store != nil {
		return store
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	store = payee.PostgresPayeeDB(db)

	return store

}

func close(store *payee.PayeePostgresDB) {
	err := store.Db.Close()
	if err != nil {
		log.Println("failed to close DB")
	}

}

func main() {
	store := initStore()
	r := payee.SetupRouter(store)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
	close(store)
}
