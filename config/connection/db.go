package connection

import (
	"database/sql"
	"sanjay/config"

	_ "github.com/lib/pq"
)

var (
	DB  *sql.DB
	err error
)

func NewDbConnection() (*sql.DB, error) {
	dbString := config.LoadEnv().GetDBURL()
	if DB, err = sql.Open("postgres", dbString); err != nil {
		return nil, err
	}
	err = DB.Ping()
	return DB, err
}
