package sales

import (
	"github.com/jinzhu/gorm"
	"github.com/strusty/worldbit-wifi/database"
)

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) database.SalesStore {
	db.AutoMigrate(&database.UsedSale{})
	return store{
		db: db,
	}
}

func (store store) Create(sale *database.UsedSale) error {
	return store.db.Create(sale).Error
}

func (store store) ByPayPalSaleID(saleID string) (*database.UsedSale, error) {
	sale := new(database.UsedSale)
	return sale, store.db.Where("pay_pal_sale_id = ?", saleID).First(sale).Error
}
