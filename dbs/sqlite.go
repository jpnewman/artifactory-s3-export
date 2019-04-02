package dbs

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/samonzeweb/godb"
	"github.com/samonzeweb/godb/adapters/sqlite"
)

// InitSqliteDb - Init SQLite DB Connection
func InitSqliteDb(filename string) *godb.DB {
	db, err := godb.Open(sqlite.Adapter, filename)
	if err != nil {
		panic(err.Error())
	}

	// OPTIONAL: Set logger to show SQL execution logs
	// db.SetLogger(log.New(os.Stderr, "", 0))

	// OPTIONAL: Set default table name building style from struct's name(if active struct doesn't have TableName() method)
	// db.SetDefaultTableNamer(tablenamer.Plural())

	return db
}

// CreateTable - Create Table
func CreateTable(db *godb.DB, obj interface{}) {
	buffer := bytes.NewBufferString("")

	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		// vField := v.Field(i)
		tField := t.Field(i)

		if i == 0 {
			buffer.WriteString("CREATE TABLE IF NOT EXISTS ")
			buffer.WriteString(t.Name())
			buffer.WriteString(" (")
		} else {
			buffer.WriteString(", ")
		}

		columnName := strings.Split(tField.Tag.Get("db"), ",")[0]
		buffer.WriteString(columnName)
		buffer.WriteString(" ")

		columnType := tField.Tag.Get("sqlite")
		buffer.WriteString(columnType)
	}

	buffer.WriteString(");")
	// fmt.Printf(buffer.String())

	_, err := db.CurrentDB().Exec(buffer.String())
	if err != nil {
		panic(fmt.Errorf("Create Table Error: %s", err))
	}
}
