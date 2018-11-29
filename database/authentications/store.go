package authentications

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/strusty/worldbit-wifi/database"
)

type store struct {
	db *gorm.DB
}

func NewAuthenticationsStore(db *gorm.DB) database.AuthenticationsStore {
	db.AutoMigrate(&database.Authentication{})
	return store{
		db: db,
	}
}

func (store store) Create(entity *database.Authentication) error {
	guid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	entity.ID = guid.String()
	return store.db.Create(entity).Error
}

func (store store) ByPhoneNumber(phoneNumber string) (*database.Authentication, error) {
	entity := new(database.Authentication)

	if err := store.db.
		Order("created_at desc").
		Where("phone_number = ?", phoneNumber).
		First(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (store store) ByConfirmationCode(confirmationCode string) (*database.Authentication, error) {
	entity := new(database.Authentication)

	if err := store.db.
		Order("created_at desc").
		Where("confirmation_code = ?", confirmationCode).
		First(entity).Error; err != nil {
		return nil, err
	}

	return entity, nil
}
