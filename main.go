package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	payee "github.com/Swarathmica-infraspec/payout-management-system/payee"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func initStore() (payee.PayeeRepository, *sql.DB) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	repo := payee.PayeeDB(db)
	return repo, db
}

func main() {
	store, db := initStore()
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("Failed to close DB:", err)
		}
	}()

	mux := payee.SetupRouter(store)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
