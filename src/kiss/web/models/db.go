package models

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Model struct {
	db   *sql.DB
	Name string
}

func New(dataSourceName string) *Model {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	return &Model{Name: "hotpie", db: db}
}

type IModel interface {
	NewEntry(url string) (int, error)
}

func (m *Model) NewEntry(url string) (int, error) {
	return 0, nil
}
