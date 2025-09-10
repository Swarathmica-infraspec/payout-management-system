package main

import (
	"database/sql"
	"log"
	"os"
	expense "payoutmanagementsystem/expense"
	payee "payoutmanagementsystem/payee"

	"github.com/gin-gonic/gin"
)

var payee_store *payee.PayeePostgresDB

var expense_store *expense.ExpensePostgresDB

func initStore() (*payee.PayeePostgresDB, *expense.ExpensePostgresDB) {
	if payee_store != nil && expense_store != nil {
		return payee_store, expense_store
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

	return payee_store, expense_store

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
	payee_store, expense_store := initStore()
	defer close(payee_store, expense_store)

	r := gin.Default()

	payeeGroup := r.Group("/")
	payee.SetupRouter(payeeGroup, payee_store)

	expenseGroup := r.Group("/")
	expense.SetupRouter(expenseGroup, expense_store)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
