package main

import (
	"database/sql"

	"github.com/TeaMeow/KitSvc/service/model"
)

func databaseUpstream(db *sql.DB) {
	_, err := db.Exec(model.TestCreateQuery)
	if err != nil {
		panic(err)
	}
}

func databaseDownstream(db *sql.DB) {

}
