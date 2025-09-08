package main

import (
	"log"
	payee "payoutmanagementsystem/payee"
)

func main() {
	r := payee.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
