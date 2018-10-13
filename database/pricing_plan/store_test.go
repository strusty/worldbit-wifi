package pricing_plan

import (
	"testing"

	"git.sfxdx.ru/crystalline/wi-fi-backend/database"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestPricingPlanStore(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")
	db.DropTableIfExists(&database.PricingPlan{})
	if assert.NoError(t, err) {
		testStore := store{
			db: db,
		}

		store := NewPricingPlanStore(db)

		if assert.Equal(t, testStore, store) {
			pricingPlan1 := &database.PricingPlan{
				ID:        "id1",
				AmountUSD: 1,
				Duration:  1,
			}
			pricingPlan2 := &database.PricingPlan{
				ID:        "id2",
				AmountUSD: 2,
				Duration:  2,
			}

			assert.NoError(t, store.Create(pricingPlan1))
			assert.NoError(t, store.Create(pricingPlan2))

			assert.NoError(t, store.Update(&database.PricingPlan{
				ID:        pricingPlan2.ID,
				AmountUSD: 3,
				Duration:  4,
			}))

			entity, err := store.ByID(pricingPlan2.ID)
			assert.NoError(t, err)
			assert.Equal(t, float64(3), entity.AmountUSD)
			assert.Equal(t, int64(4), entity.Duration)

			id := entity.ID
			assert.NoError(t, store.Delete(entity))
			entity, err = store.ByID(id)
			assert.Error(t, err)

			if allEntities, err := store.All(); assert.NoError(t, err) {
				assert.Equal(t, 1, len(allEntities))
			}
		}
	}
}
