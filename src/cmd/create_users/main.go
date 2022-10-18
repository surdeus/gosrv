package main

import (
	"encoding/json"
	"fmt"
)

type Users map[string] string
func main() {
	users := Users{
		"surdeus@gmail.com": "Password1",
		"jienfak@yandex.ru": "Password2",
	}

	usersJson, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(usersJson))
}

