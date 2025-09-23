package main

import (
	"fmt"
	"log"
	"net/http"
	payee "payoutmanagementsystem/payee"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /payees", payee.PayeePostAPI)
	fmt.Println("Server starting on :8080")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
