package models

import (
	"database/sql"
	//_ "github.com/lib/pq"
	//"log"
)

type Model struct {
	db *sql.DB

	Name   string
	lastId int
}

func New(dataSourceName string) *Model {
	//	db, err := sql.Open("postgres", dataSourceName)
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//	if err = db.Ping(); err != nil {
	//		log.Panic(err)
	//	}
	return &Model{Name: "hotpie", db: nil}
}

type IModel interface {
	//
	GetLastId() int
	// encode the long url and return the short url
	SaveTiny(url, id string) (string, error)
	//
	GetDestination(id int) (string, error)
}

func (m *Model) SaveTiny(url, id string) (string, error) {
	return id, nil
}

func (m *Model) GetLastId() int {
	m.lastId++
	return m.lastId
}

func (m *Model) GetDestination(id int) (string, error) {
	return "", nil
}
