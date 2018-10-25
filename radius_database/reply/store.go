package reply

import (
	"git.sfxdx.ru/crystalline/wi-fi-backend/radius_database"
	"github.com/jinzhu/gorm"
)

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) radius_database.ReplyStore {
	return store{
		db: db,
	}
}

func (store store) Create(check *radius_database.Reply) error {
	return store.db.Create(check).Error
}

func (store store) DeleteRepliesByUsername(username string) error {
	return store.db.Where("username = ?", username).Delete(&radius_database.Reply{}).Error
}
