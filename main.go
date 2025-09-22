package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	expense "payoutmanagementsystem/expense"
)

var store *expense.ExpensePostgresDB

func initStore() *expense.ExpensePostgresDB {
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
	store = expense.NewPostgresExpenseDB(db)
	return store
}

func main() {
	store := initStore()
	mux := http.NewServeMux()

	mux.HandleFunc("/expenses", expense.ExpensePostAPI(store))
	mux.HandleFunc("/expenses/list", expense.ExpenseGetAPI(store))
	mux.HandleFunc("/expenses/", expense.ExpenseGetOneAPI(store))

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
