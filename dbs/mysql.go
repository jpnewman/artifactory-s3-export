package dbs

import (
	"database/sql"

	"github.com/spf13/viper"
)

// InitMySQLDb - Init MySQL Database connection
func InitMySQLDb(connectionString string) *sql.DB {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}

	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_connections"))
	db.SetMaxOpenConns(viper.GetInt("mysql.max_connections"))
	db.SetConnMaxLifetime(viper.GetDuration("mysql.max_lifetime"))

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return db
}
