package storage

import (
	"database/sql"
	"fmt"
)

var dbh *sql.DB

type table struct {
	name    string
	ifField string
	fields  map[string]string
}

type record map[string]string

type Set struct{
	tableName string
	pageSize int
	pageNum int
}

type Filter struct {

}

func Init(dsn string) error {
	dbh, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("Failed to init DB")
		return err
	}

	err = dbh.Ping()
	if err != nil {
		fmt.Println("Failed to talk to DB")
		return err
	}
}

func list(query string, args ...interface{}) (sql.Result, error) {
	result, err := dbh.Exec(query, args...)
	return result, err
}

func Store(record *record) (lastInsertId string, err error) {
	return "", nil
}

func NewSet(tableName string) (*Set, error){

	// check for tableName in schema
	set := &Set{tableName:tableName}

	return set, nil
}

func (set *Set) List(query string) (sql.Result, error){

}