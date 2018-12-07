package authentications

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/strusty/worldbit-wifi/database"
)

func TestAuthenticationsStore(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")
	db.DropTableIfExists(&database.Authentication{})
	if assert.NoError(t, err) {
		testStore := store{
			db: db,
		}

		store := NewAuthenticationsStore(db)

		if assert.Equal(t, testStore, store) {
			if assert.NoError(t, store.Create(&database.Authentication{
				PhoneNumber:      "123456",
				ConfirmationCode: "code1",
			})) && assert.NoError(t, store.Create(&database.Authentication{
				PhoneNumber:      "654321",
				ConfirmationCode: "code2",
			})) && assert.NoError(t, store.Create(&database.Authentication{
				PhoneNumber:      "123456",
				ConfirmationCode: "code3",
			})) {
				entity, err := store.ByPhoneNumber("123456")
				if assert.NoError(t, err) && assert.NotNil(t, entity) {
					assert.Equal(t, "123456", entity.PhoneNumber)
					assert.Equal(t, "code3", entity.ConfirmationCode)
				}

				entity, err = store.ByPhoneNumber("1234562")
				assert.Error(t, err)
				assert.Nil(t, entity)

				entity, err = store.ByConfirmationCode("code2")
				if assert.NoError(t, err) && assert.NotNil(t, entity) {
					assert.Equal(t, "654321", entity.PhoneNumber)
					assert.Equal(t, "code2", entity.ConfirmationCode)
				}

				entity, err = store.ByConfirmationCode("1234562")
				assert.Error(t, err)
				assert.Nil(t, entity)
			}
		}
	}
}
