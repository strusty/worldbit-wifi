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

		store := New(db)
		t.Run("Initialization", func(t *testing.T) {
			if assert.Equal(t, testStore, store) {
				t.Run("Create", func(t *testing.T) {
					t.Run("Success", func(t *testing.T) {
						check := radius_database.Reply{
							Username:  "username",
							Attribute: "attribute",
							Op:        "operation",
							Value:     "value",
						}

						if assert.NoError(t, store.Create(&check)) {
							assert.Equal(t, uint(1), check.ID)
						}
					})
					db.DropTableIfExists(&radius_database.Reply{})
					t.Run("Error", func(t *testing.T) {
						assert.Error(t, store.Create(&radius_database.Reply{}))
					})
				})
			}
		})
	}
}
