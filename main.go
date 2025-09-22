package main

import (
	"fmt"
	"net/http"
	"log"
	expense "payoutmanagementsystem/expense"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /expenses", expense.ExpensePostAPI)
	fmt.Println("Server starting on :8080")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
