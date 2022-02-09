package main

import (
	"fmt"
	"mandip/go-examples/mysql"
)

type Admin struct {
	id             int64
	name           string
	email          string
	remember_token string
}

func main() {
	pool, err := mysql.Connect(mysql.DatabaseCredentials{
		"admin",
		"Root@123",
		"blood_bank_api",
		"localhost",
	})

	if err != nil {
		panic(err)
	}

	rows, _ := mysql.Run(pool, "Select id,name,email,remember_token from admins")

	for rows.Next() {
		var i Admin
		rows.Scan(&i.id, &i.name, &i.email, &i.remember_token)
		fmt.Printf("value: %v\n", i)
	}

}
