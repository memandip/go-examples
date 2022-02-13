package etl

import (
	"database/sql"
	"errors"
	"fmt"
	"mandip/go-examples/config"
	"mandip/go-examples/mysql"
)

type KV struct {
	Key, Value string
}

type DBTable struct {
	Name       string
	PrimaryKey string
	Columns    []string
}

type EtlDbCreds struct {
	CurrentDbCredentials mysql.DatabaseCredentials
	TargetDbCredentials  mysql.DatabaseCredentials
}

type Etl struct {
	CurrentSchema DBTable
	TargetSchema  DBTable
	TableMapping  []KV
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
}

func write(dataChannel <-chan map[string]interface{}, dbPool *sql.DB) {
	for result := range dataChannel {
		fmt.Println(result)
	}
}

func EtlMysql(etlDbCreds EtlDbCreds, etl Etl) {

	currentDbPool, _ := config.Dbpool(etlDbCreds.CurrentDbCredentials)
	targetDbPool, _ := config.Dbpool(etlDbCreds.TargetDbCredentials)

	err := etl.validateEtl()

	if err != nil {
		panic(err)
	}

	tableMapping := make([]KV, len(etl.CurrentSchema.Columns))

	for i, c := range etl.CurrentSchema.Columns {
		tableMapping[i] = KV{c, etl.TargetSchema.Columns[i]}
	}

	query := mysql.GenerateSelectQuery(etl.CurrentSchema.Columns, etl.CurrentSchema.Name, map[string]interface{}{})

	dataChannel := make(chan map[string]interface{})
	done := make(chan bool)
	go push(dataChannel, currentDbPool, query)
	go write(dataChannel, targetDbPool)
	<-done

}
