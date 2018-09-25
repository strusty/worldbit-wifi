package main

import (
	"fmt"
	"git.sfxdx.ru/crystalline/wi-fi-backend/database/authentications"
	"git.sfxdx.ru/crystalline/wi-fi-backend/server"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/auth"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/captcha"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/twilio"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s",
			config.GetString("postgres.host"),
			config.GetString("postgres.port"),
			config.GetString("postgres.user"),
			config.GetString("postgres.dbName"),
			config.GetString("postgres.password"),
			config.GetString("postgres.sslMode"),
			config.GetString("postgres.sslCertificatePath"),
			config.GetString("postgres.sslKeyPath"),
			config.GetString("postgres.sslRootCertificatePath"),
		),
	)
	if err != nil {
		log.Fatalf("Fatal error connecting to database: %s \n", err)
	}

	apiServer := server.New(
		auth.New(
			authentications.NewAuthenticationsStore(db),
			config.GetInt64("confirmationCodeExpiration"),
		),
		captcha.New(
			config.GetString("captcha.secret"),
		),
		twilio.New(
			config.GetString("twilio.host"),
			config.GetString("twilio.sid"),
			config.GetString("twilio.token"),
			config.GetString("twilio.from"),
			config.GetString("twilio.messageTemplate"),
		),
	)

	apiServer.Logger.Fatal(apiServer.Start(":" + config.GetString("api.port")))
}
