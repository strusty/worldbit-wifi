package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sevlyar/go-daemon"
	config "github.com/spf13/viper"
	"github.com/strusty/worldbit-wifi/database"
	"github.com/strusty/worldbit-wifi/database/admin"
	"github.com/strusty/worldbit-wifi/database/authentications"
	"github.com/strusty/worldbit-wifi/database/pricing_plan"
	"github.com/strusty/worldbit-wifi/database/sales"
	"github.com/strusty/worldbit-wifi/jwt"
	"github.com/strusty/worldbit-wifi/radius_database/accounting"
	"github.com/strusty/worldbit-wifi/radius_database/check"
	"github.com/strusty/worldbit-wifi/radius_database/reply"
	"github.com/strusty/worldbit-wifi/server"
	"github.com/strusty/worldbit-wifi/services/admins"
	"github.com/strusty/worldbit-wifi/services/auth"
	"github.com/strusty/worldbit-wifi/services/captcha"
	"github.com/strusty/worldbit-wifi/services/cleaner"
	"github.com/strusty/worldbit-wifi/services/paypal"
	"github.com/strusty/worldbit-wifi/services/pricing_plans"
	"github.com/strusty/worldbit-wifi/services/radius"
	"github.com/strusty/worldbit-wifi/services/twilio"
	"github.com/strusty/worldbit-wifi/services/worldbit"
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

	if config.GetBool("api.daemon") {
		log.Println("Create daemon context")
		context := &daemon.Context{
			PidFileName: "wi_fi_backend_pidfile",
			PidFilePerm: 0644,
			LogFileName: "wi_fi_backend_logfile",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
		}
		log.Println("Context rebirth")

		d, err := context.Reborn()
		if err != nil {
			log.Fatal("Unable to run: " + err.Error())
		}

		if d != nil {
			log.Printf("Daemon is running, and its pid is: %d \n", d.Pid)
			return
		}
		defer log.Println("release context")
		defer context.Release()

		log.Println("- - - - - - - - - - - - - - -")
		log.Println("daemon started")
	}

	apiServer.Logger.Fatal(apiServer.Start(":" + config.GetString("api.port")))
}
