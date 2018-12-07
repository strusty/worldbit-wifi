package reply

import (
	"github.com/jinzhu/gorm"
	"github.com/strusty/worldbit-wifi/radius_database"
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
