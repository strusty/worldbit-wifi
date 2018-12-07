package check

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/strusty/worldbit-wifi/radius_database"
)

func TestRadiusCheckStore(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")
	db.AutoMigrate(&radius_database.Check{})
	if assert.NoError(t, err) {
		testStore := store{
			db: db,
		}

		store := NewStore(db)
		if assert.Equal(t, testStore, store) {
			check := radius_database.Check{
				Username:  "username",
				Attribute: "Max-Daily-Session",
				Op:        "operation",
				Value:     "1234",
			}

			if assert.NoError(t, store.Create(&check)) {
				assert.Equal(t, uint(1), check.ID)
			}

			check = radius_database.Check{
				Username:  "usernames",
				Attribute: "Max-Daily-Session",
				Op:        "operation",
				Value:     "12345",
			}

			if assert.NoError(t, store.Create(&check)) {
				assert.Equal(t, uint(2), check.ID)
			}

			check = radius_database.Check{
				Username:  "user",
				Attribute: "Max-Daily-Session",
				Op:        "operation",
				Value:     "12345",
			}

			if assert.NoError(t, store.Create(&check)) {
				assert.Equal(t, uint(3), check.ID)
			}

			check = radius_database.Check{
				Username:  "user",
				Attribute: "attribute",
				Op:        "operation",
				Value:     "12345",
			}

			if assert.NoError(t, store.Create(&check)) {
				assert.Equal(t, uint(4), check.ID)
			}

			checks, err := store.SessionChecks()
			if assert.NoError(t, err) && assert.Len(t, checks, 3) {
				assert.Equal(t, radius_database.Check{
					ID:        3,
					Username:  "user",
					Attribute: "Max-Daily-Session",
					Op:        "operation",
					Value:     "12345",
				}, checks[2])
			}

			if assert.NoError(t, store.DeleteChecksByUsername("user")) {
				checks, err := store.SessionChecks()
				assert.NoError(t, err)
				assert.Len(t, checks, 2)
			}

			db.DropTableIfExists(&radius_database.Check{})
			assert.Error(t, store.Create(&radius_database.Check{}))
			_, err = store.SessionChecks()
			assert.Error(t, err)
		}
	}
}
