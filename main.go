package main

import (
	"database/sql"
	"log"
	"os"
	payee "payoutmanagementsystem/payee"
	expense "payoutmanagementsystem/expense"
)

var payee_store *payee.PayeePostgresDB


var expense_store *expense.ExpensePostgresDB

func initStore() (*payee.PayeePostgresDB,*expense.ExpensePostgresDB) {
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
	payee_store = payee.PostgresPayeeDB(db)
	expense_store = expense.NewPostgresExpenseDB(db)

	return payee_store,expense_store

}

func close(payee_store *payee.PayeePostgresDB, expense_store *expense.ExpensePostgresDB) {
	err := payee_store.Db.Close()
	if err != nil {
		log.Println("failed to close DB")
	}
	err = expense_store.Db.Close()
	if err != nil {
		log.Println("failed to close DB")
	}
}

func main() {
	payee_store,expense_store := initStore()
	r := payee.SetupRouter(payee_store)
	r1 := expense.SetupRouter(expense_store)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

	if err := r1.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
	close(payee_store,expense_store)
}
