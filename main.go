package main

import (
	"database/sql"
	"log"
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

func close(store *expense.ExpensePostgresDB) {
	err := store.Db.Close()
	if err != nil {
		log.Println("failed to close DB")
	}

}

func main() {
	store := initStore()
	r := expense.SetupRouter(store)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
	close(store)
}
