package connection

import (
	"database/sql"
	"fmt"
	"net/url"
	"sanjay/config"

	_ "github.com/lib/pq"
)

var (
	DB  *sql.DB
	err error
)

func NewDbConnection() (*sql.DB, error) {
	dbString := config.LoadEnv().GetDBURL()
	fmt.Println("str--------->", dbString)
	conn, _ := url.Parse(dbString)
	conn.RawQuery = "sslmode=verify-ca;sslrootcert=ca.pem"
	if DB, err = sql.Open("postgres", conn.String()); err != nil {
		return nil, err
	}
	err = DB.Ping()
	return DB, err
}
