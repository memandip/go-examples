package main

import (
	"fmt"
	"mandip/go-examples/config"
	"mandip/go-examples/mysql"
)

func main() {
	pool, err := config.Dbpool()

	if err != nil {
		panic(err)
	}

	selection := []string{"id", "name", "email"}
	query := mysql.GenerateSelectQuery(selection, "admins", map[string]interface{}{
		"id":    1,
		"email": "superadmin@gmail.com",
		"name":  "Superadmin",
	})
	fmt.Println(query)
	rows, _ := mysql.Run(pool, query)

	for rows.Next() {
		var i mysql.Admin
		rows.Scan(&i.Id, &i.Name, &i.Email)
		fmt.Printf("value: %v\n", i)
	}

}
