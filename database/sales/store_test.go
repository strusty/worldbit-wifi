package sales

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/strusty/worldbit-wifi/database"
)

func TestPayPalStore(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")
	db.DropTableIfExists(&database.UsedSale{})
	if assert.NoError(t, err) {
		testStore := store{
			db: db,
		}

		store := NewStore(db)

		if assert.Equal(t, testStore, store) {
			if assert.NoError(t, store.Create(&database.UsedSale{
				PayPalSaleID: "id",
				Voucher:      "voucher",
			})) && assert.NoError(t, store.Create(&database.UsedSale{
				PayPalSaleID: "id1",
				Voucher:      "voucher1",
			})) && assert.NoError(t, store.Create(&database.UsedSale{
				PayPalSaleID: "id2",
				Voucher:      "voucher2",
			})) {
				sale, err := store.ByPayPalSaleID("id1")
				if assert.NoError(t, err) {
					assert.Equal(t, database.UsedSale{
						PayPalSaleID: "id1",
						Voucher:      "voucher1",
					}, *sale)
				}

				_, err = store.ByPayPalSaleID("id3")
				assert.Error(t, err)
			}
		}
	}
}
