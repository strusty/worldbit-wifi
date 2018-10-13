package admins

import (
	"git.sfxdx.ru/crystalline/wi-fi-backend/database"
	"git.sfxdx.ru/crystalline/wi-fi-backend/jwt"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	store database.AdminStore
}

func New(store database.AdminStore) Admins {
	return service{
		store: store,
	}
}

func (service service) Login(request LoginRequest) (*JWTResponse, error) {
	admin, err := service.store.ByLogin(request.Login)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(request.Password)); err != nil {
		return nil, err
	}

	token, err := jwt.GenerateJWT(admin.ID)
	if err != nil {
		return nil, err
	}

	return &JWTResponse{
		Token: token,
	}, nil
}

func (service service) ChangePassword(adminID string, request ChangePasswordRequest) error {
	admin, err := service.store.ByID(adminID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(request.OldPassword)); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 15)
	if err != nil {
		return err
	}

	return service.store.Update(adminID, "password", string(hash))
}
