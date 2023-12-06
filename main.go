package main

import (
	"fmt"

	"github.com/Kotletta-TT/bonus-service/internal/utils"
)

func main() {
	user := "user1"
	secretKey := "secret1"
	token, err := utils.GenerateToken(user, secretKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(token)
}
