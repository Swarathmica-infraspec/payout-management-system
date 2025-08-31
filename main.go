package main

import (
	"log"
	payee "payoutmanagementsystem/payee"
)

func Print(s string) string {
	return s
}

func main() {
	r := payee.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
