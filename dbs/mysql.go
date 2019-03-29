package dbs

import "database/sql"

func InitMySqlDb(connectionString string) *sql.DB {
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return db
}
