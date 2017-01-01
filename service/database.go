package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// createDatabase creates the database connection.
func createDatabase() *gorm.DB {

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.name"),
		viper.GetString("database.charset"),
		viper.GetString("database.parse_time"),
		viper.GetString("database.loc"),
	))

	if err != nil {
		panic(err)
	}

	defer db.Close()

	return db
}
