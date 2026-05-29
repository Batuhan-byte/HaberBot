package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash, err := bcrypt.GenerateFromPassword([]byte("user"), 10)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Hash:", string(hash))
}
