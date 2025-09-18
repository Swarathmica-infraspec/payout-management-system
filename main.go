package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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
	mux := http.NewServeMux()
	mux.HandleFunc("POST /payee", payee.PayeePostAPI)
	mux.HandleFunc("GET /payee", payee.PayeeGetAPI(store))
	mux.HandleFunc("GET /payee/:id", payee.PayeeGetOneAPI(store))
	fmt.Println("Server starting on :8080")
	mux.HandleFunc("/payee/", payee.PayeeGetOneAPI(store))

	log.Fatal(http.ListenAndServe(":8080", mux))
	close(store)
}
