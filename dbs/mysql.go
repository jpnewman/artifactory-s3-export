package dbs

import (
	"database/sql"

	"github.com/spf13/viper"
)

func InitMySqlDb(connectionString string) *sql.DB {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}

	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(viper.GetInt("mysql.max_connections"))
	db.SetConnMaxLifetime(viper.GetDuration("mysql.max_lifetime"))

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return db
}
