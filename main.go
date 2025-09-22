package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	payee "payoutmanagementsystem/payee"

	_ "github.com/lib/pq"
)

var store payee.PayeeRepository

func initStore() payee.PayeeRepository {
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
	store = payee.PayeeDB(db)
	return store
}

func main() {
	store := initStore()
	mux := http.NewServeMux()

	mux.HandleFunc("/payees", payee.PayeePostAPI(store))
	mux.HandleFunc("/payees/list", payee.PayeeGetAPI(store))
	mux.HandleFunc("/payees/", payee.PayeeGetOneAPI(store))
	mux.HandleFunc("/payees/update/", payee.PayeeUpdateAPI(store))

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
