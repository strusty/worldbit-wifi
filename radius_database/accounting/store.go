package accounting

import (
	"github.com/jinzhu/gorm"
	"github.com/strusty/worldbit-wifi/radius_database"
)

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) radius_database.AccountingStore {
	return store{
		db: db,
	}
}

func (store store) SessionTimeSum(username string) (int64, error) {
	type sqlsum struct {
		Sum int64
	}

	sum := sqlsum{}
	return sum.Sum, store.db.
		Raw("SELECT SUM(AcctSessionTime) as sum FROM radacct WHERE username=?", username).
		Scan(&sum).Error
}

func (store store) DeleteAccountingByUsername(username string) error {
	return store.db.Exec("DELETE FROM radacct WHERE username = ?", username).Error
}
