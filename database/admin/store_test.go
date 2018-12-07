package admin

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/strusty/worldbit-wifi/database"
)

func TestAdminStore(t *testing.T) {
	db, err := gorm.Open("sqlite3", "test.db")
	db.DropTableIfExists(&database.Admin{})

	if assert.NoError(t, err) {
		testStore := store{
			db: db,
		}

		store := NewAdminStore(db)

		if assert.Equal(t, testStore, store) {

			admin1 := &database.Admin{
				Login:    "adm1",
				Password: "pass1",
			}
			admin2 := &database.Admin{
				Login:    "adm2",
				Password: "pass2",
			}

			assert.NoError(t, store.Create(admin1))
			assert.NoError(t, store.Create(admin2))

			assert.NoError(t, store.Update(admin2.ID, "Password", "newpass"))

			entity, err := store.ByLogin("adm2")
			assert.NoError(t, err)
			assert.Equal(t, "adm2", entity.Login)
			assert.Equal(t, "newpass", entity.Password)

			entity, err = store.ByID(admin1.ID)
			assert.NoError(t, err)
			assert.Equal(t, "adm1", entity.Login)
			assert.Equal(t, "pass1", entity.Password)

		}
	}
}
