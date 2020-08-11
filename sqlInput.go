package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Row is one row
type Row map[string]string

// Input is the generic input interface
type Input interface {
	Execute(query string) ([]Row, error)
}

type sqlInput struct {
	conn *sql.DB
}

func (i *sqlInput) Execute(query string) (rows []Row, err error) {
	return
}

// MakeSQLInput creates an input and returns that interface
func MakeSQLInput(conf InputConfig) Input {
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

	return &sqlInput{conn: db}
}
