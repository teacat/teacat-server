package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// createDatabase creates the database connection.
func createDatabase(resetDB *bool) *gorm.DB {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%s&loc=%s",
		os.Getenv("KITSVC_DATABASE_USER"),
		os.Getenv("KITSVC_DATABASE_PASSWORD"),
		os.Getenv("KITSVC_DATABASE_HOST"),
		os.Getenv("KITSVC_DATABASE_NAME"),
		os.Getenv("KITSVC_DATABASE_CHARSET"),
		os.Getenv("KITSVC_DATABASE_PARSE_TIME"),
		os.Getenv("KITSVC_DATABASE_LOC"),
	))

	if err != nil {
		panic(err)
	}

	defer db.Close()

	return db
}
