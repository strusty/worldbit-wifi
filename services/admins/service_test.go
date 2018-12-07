package admins

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/strusty/worldbit-wifi/database"
	"github.com/strusty/worldbit-wifi/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AdminStoreMock struct {
	CreateFn  func(admin *database.Admin) error
	UpdateFn  func(id string, key string, value interface{}) error
	ByLoginFn func(login string) (*database.Admin, error)
	ByIDFn    func(id string) (*database.Admin, error)
}

func (mock AdminStoreMock) Create(entity *database.Admin) error {
	return mock.CreateFn(entity)
}

func (mock AdminStoreMock) Update(id string, key string, value interface{}) error {
	return mock.UpdateFn(id, key, value)
}

func (mock AdminStoreMock) ByLogin(login string) (*database.Admin, error) {
	return mock.ByLoginFn(login)
}

func (mock AdminStoreMock) ByID(id string) (*database.Admin, error) {
	return mock.ByIDFn(id)
}

func TestNew(t *testing.T) {
	store := AdminStoreMock{}

	testService := service{
		store: &store,
	}

	service := New(&store)

	assert.Equal(t, testService, service)
}

func CreateFn(admin *database.Admin) error {
	admin.ID = "id"
	return nil
}

func Test_service_Login(t *testing.T) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte("password"), 15)
	assert.NoError(t, err)
	jwt.SetRandomSecret()

	admin := &database.Admin{
		Login:    "admin",
		Password: string(hashedPass),
	}

	CreateFn(admin)

	service := service{
		store: AdminStoreMock{
			ByLoginFn: func(login string) (*database.Admin, error) {
				if login == "admin" {
					return admin, nil
				} else {
					return nil, errors.New("No such a user")
				}
			},
		},
	}

	t.Run("Succesful login", func(t *testing.T) {
		jwtResponse, err := service.Login(LoginRequest{
			Login:    "admin",
			Password: "password",
		})
		assert.NoError(t, err)
		assert.NotNil(t, jwtResponse)
	})

	t.Run("Login failed: wrong password", func(t *testing.T) {
		jwtResponse, err := service.Login(LoginRequest{
			Login:    "admin",
			Password: "buzzword",
		})
		assert.Error(t, err)
		assert.Nil(t, jwtResponse)
	})

	t.Run("Login failed: wrong login", func(t *testing.T) {
		jwtResponse, err := service.Login(LoginRequest{
			Login:    "odmen",
			Password: "password",
		})
		assert.Error(t, err)
		assert.Nil(t, jwtResponse)
	})
}

func Test_service_ChangePassword(t *testing.T) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte("password"), 15)
	assert.NoError(t, err)
	jwt.SetRandomSecret()

	admin := &database.Admin{
		Login:    "admin",
		Password: string(hashedPass),
	}

	CreateFn(admin)

	service := service{
		store: AdminStoreMock{
			ByIDFn: func(id string) (*database.Admin, error) {
				if id == "id" {
					return admin, nil
				} else {
					return nil, errors.New("No such a user")
				}
			},
			UpdateFn: func(id string, key string, value interface{}) error {
				if id == "id" && key == "password" {
					admin.Password = value.(string)
				}
				return nil
			},
		},
	}

	t.Run("Password change success", func(t *testing.T) {
		err := service.ChangePassword("id", ChangePasswordRequest{
			OldPassword: "password",
			NewPassword: "newpass",
		})
		assert.NoError(t, err)
	})

	t.Run("Password change fail: wrong id", func(t *testing.T) {
		err := service.ChangePassword("wrong_id", ChangePasswordRequest{
			OldPassword: "password",
			NewPassword: "newpass",
		})
		assert.Error(t, err)
	})

	t.Run("Password change fail: wrong pass", func(t *testing.T) {
		err := service.ChangePassword("id", ChangePasswordRequest{
			OldPassword: "wrong_pass",
			NewPassword: "newpass",
		})
		assert.Error(t, err)
	})
}
