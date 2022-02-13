package mysql

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func getDBCredentials() DatabaseCredentials {
	godotenv.Load("../.env.test")

	return DatabaseCredentials{
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_HOST"),
	}
}

func TestConnect(t *testing.T) {

	tests := []struct {
		name        string
		credentials DatabaseCredentials
	}{
		{"Connect", getDBCredentials()},
		{"Connect", getDBCredentials()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, err := Connect(tt.credentials)

			if pool.Ping() != nil {
				t.Errorf("Connect() error = %v", err)
				return
			}

		})
	}
}

func TestRun(t *testing.T) {
	dbPool, _ := Connect(getDBCredentials())
	tests := []struct {
		name    string
		dbPool  *sql.DB
		query   string
		wantErr bool
	}{
		{"Run", dbPool, "SELECT id, name, email FROM admins", false},
		{"Run", dbPool, "SELECT *, FROM admins", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := Run(tt.dbPool, tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			defer rows.Close()
			for rows.Next() {
				var id int
				var name, email string
				err := rows.Scan(&id, &name, &email)
				if err != nil {
					t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}
