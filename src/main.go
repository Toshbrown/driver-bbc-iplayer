package main

import (
	"fmt"
)

func main() {

	fmt.Println("get token")
	token, err := Auth("someuser@gmail.com", "somePassword")

	fmt.Println(token, err)

	recommendations, err := GetRecommendations(token)

	fmt.Println(recommendations, err)
}
