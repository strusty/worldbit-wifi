package main

import (
	"fmt"
	"log"
	"time"

	"git.sfxdx.ru/crystalline/wi-fi-backend/database"
	"git.sfxdx.ru/crystalline/wi-fi-backend/database/admin"
	"git.sfxdx.ru/crystalline/wi-fi-backend/database/authentications"
	"git.sfxdx.ru/crystalline/wi-fi-backend/database/pricing_plan"
	"git.sfxdx.ru/crystalline/wi-fi-backend/database/sales"
	"git.sfxdx.ru/crystalline/wi-fi-backend/jwt"
	"git.sfxdx.ru/crystalline/wi-fi-backend/radius_database/accounting"
	"git.sfxdx.ru/crystalline/wi-fi-backend/radius_database/check"
	"git.sfxdx.ru/crystalline/wi-fi-backend/radius_database/reply"
	"git.sfxdx.ru/crystalline/wi-fi-backend/server"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/admins"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/auth"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/captcha"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/cleaner"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/paypal"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/pricing_plans"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/radius"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/twilio"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/worldbit"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	config "github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Println("GET CONFIG")

	config.AddConfigPath("./")
	config.SetConfigName("config")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error getting config from file: %s \n", err)
	}

	// config := config.Sub(prefix)

	if err := jwt.SetRandomSecret(); err != nil {
		log.Fatalf("Fatal error setting jwt secret: %s \n", err)
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

	radiusDB, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s",
			config.GetString("radusDatabase.host"),
			config.GetString("radusDatabase.port"),
			config.GetString("radusDatabase.user"),
			config.GetString("radusDatabase.dbName"),
			config.GetString("radusDatabase.password"),
			config.GetString("radusDatabase.sslMode"),
			config.GetString("radusDatabase.sslCertificatePath"),
			config.GetString("radusDatabase.sslKeyPath"),
			config.GetString("radusDatabase.sslRootCertificatePath"),
		),
	)
	if err != nil {
		log.Fatalf("Fatal error connecting to database: %s \n", err)
	}

	adminStore := admin.NewAdminStore(db)
	adminLoginDefault := config.GetString("admin.login")
	adminPasswordDefault := config.GetString("admin.password")
	_, err = adminStore.ByLogin(adminLoginDefault)
	if err != nil {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(adminPasswordDefault), 15)
		if err != nil {
			log.Fatalf("Unable to pregenerate admin: %s\n", err)
		}

		if err := adminStore.Create(&database.Admin{
			ID:       "admin_id",
			Login:    adminLoginDefault,
			Password: string(hashedPass),
		}); err != nil {
			log.Fatalf("Unable to pregenerate admin: %s\n", err)
		}
	}

	accountingStore := accounting.NewStore(radiusDB)
	checkStore := check.NewStore(radiusDB)
	replyStore := reply.NewStore(radiusDB)

	cleanerService := cleaner.New(
		accountingStore,
		checkStore,
		replyStore,
	)

	go cleanerService.Start(time.Minute)

	apiServer, err := server.New(
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
			config.GetString("twilio.confirmationCodeMessageTemplate"),
			config.GetString("twilio.voucherMessageTemplate"),
		),
		worldbit.New(
			worldbit.Config{
				APIKey:            config.GetString("worldbit.apiKey"),
				APISecret:         config.GetString("worldbit.apiSecret"),
				MerchantID:        config.GetString("worldbit.merchantID"),
				Host:              config.GetString("worldbit.host"),
				MonitoringTimeout: config.GetInt64("worldbit.monitoringTimeout"),
				DefaultCurrency:   config.GetString("worldbit.defaultCurrency"),
				DefaultEmail:      config.GetString("worldbit.defaultEmail"),
			},
		),
		radius.New(
			checkStore,
			replyStore,
		),
		pricing_plans.New(
			pricing_plan.NewPricingPlanStore(db),
		),
		admins.New(
			adminStore,
		),
		paypal.New(
			sales.NewStore(db),
			paypal.Config{
				Host:     config.GetString("paypal.host"),
				ClientID: config.GetString("paypal.clientID"),
				Secret:   config.GetString("paypal.secret"),
			},
		),
	)
	if err != nil {
		log.Fatalf("Unable to init api server: %s\n", err)
	}

	apiServer.Logger.Fatal(apiServer.Start(":" + config.GetString("api.port")))
}
