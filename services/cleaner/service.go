package cleaner

import (
	"log"
	"strconv"
	"time"

	"github.com/strusty/worldbit-wifi/radius_database"
)

type service struct {
	accountingStore radius_database.AccountingStore
	checkStore      radius_database.CheckStore
	replyStore      radius_database.ReplyStore
}

func New(
	accountingStore radius_database.AccountingStore,
	checkStore radius_database.CheckStore,
	replyStore radius_database.ReplyStore,
) Cleaner {
	return service{
		accountingStore: accountingStore,
		checkStore:      checkStore,
		replyStore:      replyStore,
	}
}

func (service service) Start(period time.Duration) {
	tickerChannel := time.Tick(period)

	for range tickerChannel {
		service.cleanup()
	}
}

func (service service) cleanup() {
	checks, err := service.checkStore.SessionChecks()
	if err != nil {
		log.Printf("Unable to retrieve checks from database: %s\n", err)
		return
	}

	for _, check := range checks {
		sessionTime, err := strconv.ParseInt(check.Value, 10, 64)
		if err != nil {
			log.Printf("Unable to parse attribute %s:%s for username %s:%s\n", check.Attribute, check.Value, check.Username, err)
			continue
		}

		timeSum, err := service.accountingStore.SessionTimeSum(check.Username)
		if err != nil {
			log.Printf("Unable to session time sum for user %s from database: %s\n", check.Username, err)
		}

		if timeSum >= sessionTime {
			if err := service.checkStore.DeleteChecksByUsername(check.Username); err != nil {
				log.Printf("Unable to delete checks for user %s", check.Username)
				continue
			}
			if err := service.replyStore.DeleteRepliesByUsername(check.Username); err != nil {
				log.Printf("Unable to delete replies for user %s", check.Username)

			}
			if err := service.accountingStore.DeleteAccountingByUsername(check.Username); err != nil {
				log.Printf("Unable to delete accounting entries for user %s", check.Username)
			}
		}
	}
}
