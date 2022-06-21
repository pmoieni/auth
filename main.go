package main

import (
	api "github.com/pmoieni/auth/api/v1"
	"github.com/pmoieni/auth/config"
	"github.com/pmoieni/auth/log"
	"github.com/pmoieni/auth/store"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Logger.Fatal("failed to read env file: " + err.Error())
	}

	dbConfig := store.DBConfig{
		DBHost:     config.DBHost,
		DBUser:     config.DBUser,
		DBPassword: config.DBPassword,
		DBName:     config.DBName,
		DBPort:     config.DBPort,
	}

	err = store.InitDB(&dbConfig)
	if err != nil {
		log.Logger.Fatal("failed to establish connection to database: " + err.Error())
	}

	err = store.CreateRefreshTokenDatabase("localhost:6379", "")
	if err != nil {
		log.Logger.Fatal("failed to connect to redis client: " + err.Error())
	}

	err = store.CreateClientIDDatabase("localhost:6379", "")
	if err != nil {
		log.Logger.Fatal("failed to connect to redis client: " + err.Error())
	}

	err = store.CreatePasswordResetTokenDatabase("localhost:6379", "")
	if err != nil {
		log.Logger.Fatal("failed to connect to redis client: " + err.Error())
	}

	server := api.Server{}
	server.Run(":8080")
}
