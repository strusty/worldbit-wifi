package reply

import (
	"testing"

	"git.sfxdx.ru/crystalline/wi-fi-backend/radius_database"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestRadiusCheckStore(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")
	db.AutoMigrate(&radius_database.Reply{})
	if assert.NoError(t, err) {
		testStore := store{
			db: db,
		}

		store := NewStore(db)
		if assert.Equal(t, testStore, store) {

			check := radius_database.Reply{
				Username:  "username",
				Attribute: "attribute",
				Op:        "operation",
				Value:     "value",
			}

			if assert.NoError(t, store.Create(&check)) {
				assert.Equal(t, uint(1), check.ID)
			}

			check = radius_database.Reply{
				Username:  "username",
				Attribute: "attribute",
				Op:        "operation",
				Value:     "value",
			}

			if assert.NoError(t, store.Create(&check)) {
				assert.Equal(t, uint(2), check.ID)
			}

			check = radius_database.Reply{
				Username:  "usernames",
				Attribute: "attribute",
				Op:        "operation",
				Value:     "value",
			}

			if assert.NoError(t, store.Create(&check)) {
				assert.Equal(t, uint(3), check.ID)
			}

			assert.NoError(t, store.DeleteRepliesByUsername("username"))

			checks := make([]radius_database.Reply, 0)
			if assert.NoError(t, db.Find(&checks).Error) {
				assert.Len(t, checks, 1)
			}

			db.DropTableIfExists(&radius_database.Reply{})
			assert.Error(t, store.Create(&radius_database.Reply{}))
		}
	}
}
