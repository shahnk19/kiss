package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"kiss/web/baseEnc"
	"log"
	"time"
)

const (
	SQL_CREATE_TABLE = `CREATE TABLE IF NOT EXISTS tinyurlmap
    (
        id serial PRIMARY KEY NOT NULL,
        url varchar(500) NOT NULL,
        created integer NOT NULL DEFAULT 0,
        unique (url)
    );`
	SQL_GET_LAST_ID   = `SELECT id FROM tinyurlmap ORDER BY id DESC LIMIT 1;`
	SQL_GET_URL_BY_ID = `SELECT url FROM tinyurlmap WHERE id=%d;`
	SQL_GET_ID_BY_URL = `SELECT id FROM tinyurlmap WHERE url='%s';`
	SQL_INSERT        = `INSERT INTO tinyurlmap 
		(url,created)
		VALUES ('%s',%d) 
		RETURNING id;`
)

type Model struct {
	db *sql.DB

	Name   string
	lastId int
}

func New(dataSourceName string) *Model {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	if _, err = db.Exec(SQL_CREATE_TABLE); err != nil {
		log.Panic(err)
	}

	return &Model{Name: "hotpie", db: db}
}

type IModel interface {
	//
	GetLastId() int
	// encode the long url and return the short url
	SaveTiny(url string, enc *baseEnc.Encoding) (string, error)
	//
	GetIdByUrl(url string) (int, error)
	GetUrlById(id int) (string, error)
}

func (m *Model) SaveTiny(url string, enc *baseEnc.Encoding) (string, error) {
	sql := fmt.Sprintf(SQL_INSERT, url, time.Now().Unix())
	var lastInsertId int
	err := m.db.QueryRow(sql).Scan(&lastInsertId)
	if err != nil {
		log.Println(err)
		return "", err
	}
	tiny := enc.BaseEncode(lastInsertId)
	m.lastId = lastInsertId
	return tiny, nil
}

func (m *Model) GetLastId() int {
	if m.lastId <= 0 {
		if err := m.getLastIdFromDb(); err != nil {
			m.lastId = 0
		}
	}
	m.lastId++
	return m.lastId
}

func (m *Model) GetUrlById(id int) (string, error) {
	var url string
	err := m.db.QueryRow(fmt.Sprintf(SQL_GET_URL_BY_ID, id)).Scan(&url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return url, nil
}

func (m *Model) GetIdByUrl(url string) (int, error) {
	var id int
	err := m.db.QueryRow(fmt.Sprintf(SQL_GET_ID_BY_URL, url)).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return id, nil
}

func (m *Model) getLastIdFromDb() error {
	var lastInsertId int
	err := m.db.QueryRow(SQL_GET_LAST_ID).Scan(&lastInsertId)
	if err != nil {
		if err == sql.ErrNoRows {
			m.lastId = 0
			return nil
		}
		log.Println(err)
		return err
	}
	m.lastId = lastInsertId
	return nil
}
