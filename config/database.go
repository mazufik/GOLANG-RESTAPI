package config

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	viper.SetConfigFile("config.json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	dbname := viper.GetString("database.dbname")
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")

	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})

	if err != nil {
		panic("Can't connect to database")
	}

	return db
}
