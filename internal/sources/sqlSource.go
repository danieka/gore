package sources

import (
	"database/sql"
	"fmt"
	"log"

	// import mysql
	_ "github.com/go-sql-driver/mysql"
)

// Row is one row
type Row map[string]interface{}

// Source is the generic input interface
type Source interface {
	Execute(query string) ([]string, []Row, error)
}

type sqlInput struct {
	conn *sql.DB
}

func (i *sqlInput) Execute(query string) (cols []string, rows []Row, err error) {
	dbRows, err := i.conn.Query(query)
	if err != nil {
		return
	}
	defer dbRows.Close()

	cols, err = dbRows.Columns()
	if err != nil {
		return
	}

	for dbRows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err = dbRows.Scan(columnPointers...); err != nil {
			return
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			f, ok := (*val).([]byte)
			if ok {
				m[colName] = string(f)
			} else {
				fmt.Println(*val)
				m[colName] = *val
			}
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		rows = append(rows, m)
	}
	return
}

// MakeSQLSource creates an input and returns that interface
func MakeSQLSource(key string, conf SourceConfig) {
	if conf.Type != "mysql" {
		log.Fatal("Database not supported " + conf.Type)
	}
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	Sources[key] = &sqlInput{conn: db}
}
