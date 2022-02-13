package etl

import (
	"database/sql"
	"errors"
	"fmt"
	"mandip/go-examples/config"
	"mandip/go-examples/mysql"
	"strings"
)

type KV struct {
	Key, Value string
}

type DBTable struct {
	Name, PrimaryKey string
	Columns          []string
}

type EtlDbCreds struct {
	CurrentDbCredentials, TargetDbCredentials mysql.DatabaseCredentials
}

type Etl struct {
	CurrentSchema, TargetSchema DBTable
	TableMapping                map[string]string
}

func (etl *Etl) validateEtl() error {
	if len(etl.CurrentSchema.Name) == 0 || len(etl.TargetSchema.Name) == 0 {
		return errors.New("current schema name and target schema name is required")
	}

	if len(etl.CurrentSchema.Columns) != len(etl.TargetSchema.Columns) {
		return errors.New("total columns does not match with current and target schema")
	}

	return nil
}

func push(dataChannel chan<- map[string]interface{}, dbPool *sql.DB, query string) {
	rows, _ := mysql.Run(dbPool, query)

	for rows.Next() {
		var id int64
		var name, email string
		rows.Scan(&id, &name, &email)
		result := map[string]interface{}{
			"id":    id,
			"name":  name,
			"email": email,
		}
		dataChannel <- result
	}

	defer rows.Close()
	defer close(dataChannel)
}

func write(dataChannel <-chan map[string]interface{}, dbPool *sql.DB, etl Etl, done chan<- bool) {
	for {
		result, more := <-dataChannel

		if !more {
			done <- true
			break
		}

		query := "INSERT INTO " + etl.TargetSchema.Name
		var columns, values []string
		for k, v := range result {
			targetKey := etl.TableMapping[k]
			if len(targetKey) == 0 {
				break
			}
			columns = append(columns, targetKey)
			values = append(values, fmt.Sprintf("'%v'", v))
		}

		query += "(`" + strings.Join(columns, "`, `") + "`) VALUES (" + strings.Join(values, ", ") + ")"
		_, err := mysql.Run(dbPool, query)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func EtlMysql(etlDbCreds EtlDbCreds, etl Etl) {

	currentDbPool, _ := config.Dbpool(etlDbCreds.CurrentDbCredentials)
	targetDbPool, _ := config.Dbpool(etlDbCreds.TargetDbCredentials)

	if currentDbPool.Ping() != nil {
		fmt.Println("current db not connected")
	}

	if targetDbPool.Ping() != nil {
		fmt.Println("target db not connected")
	}

	err := etl.validateEtl()

	if err != nil {
		panic(err)
	}

	tableMapping := make(map[string]string, len(etl.CurrentSchema.Columns))

	for i, c := range etl.CurrentSchema.Columns {
		tableMapping[c] = etl.TargetSchema.Columns[i]
	}

	etl.TableMapping = tableMapping

	query := mysql.GenerateSelectQuery(etl.CurrentSchema.Columns, etl.CurrentSchema.Name, map[string]interface{}{})

	dataChannel := make(chan map[string]interface{})
	done := make(chan bool)
	go push(dataChannel, currentDbPool, query)
	go write(dataChannel, targetDbPool, etl, done)
	<-done

}
