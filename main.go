package main

import (
	"git.sfxdx.ru/crystalline/wi-fi-backend/server"
	config "github.com/spf13/viper"
	"log"
)

func main() {
	log.Println("GET CONFIG")

	config.AddConfigPath("./")
	config.SetConfigName("config")

	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error getting config from file: %s \n", err)
	}

	apiServer := server.New()

	apiServer.Logger.Fatal(apiServer.Start(":" + config.GetString("api.port")))
}
