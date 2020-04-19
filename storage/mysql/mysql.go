package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DB struct {
	client *gorm.DB
}

// New creates database connection
func New(user, password, host, database string, port int) (*DB, error) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8", user, password, host, port, database))
	if err != nil {
		return &DB{}, err
	}

	return &DB{client: db}, nil
}

// Close closes the connection to database
func (db *DB) Close() {
	if db.client != nil {
		db.client.Close()
	}
}
