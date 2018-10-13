package pricing_plans

import (
	"testing"

	"git.sfxdx.ru/crystalline/wi-fi-backend/database"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type PricingPlanStoreMock struct {
	CreateFn func(entity *database.PricingPlan) error
	UpdateFn func(entity *database.PricingPlan) error
	DeleteFn func(id string) error
	AllFn    func() ([]database.PricingPlan, error)
	ByIDFn   func(id string) (*database.PricingPlan, error)
}

func (mock PricingPlanStoreMock) Create(entity *database.PricingPlan) error {
	return mock.CreateFn(entity)
}

func (mock PricingPlanStoreMock) Update(entity *database.PricingPlan) error {
	return mock.UpdateFn(entity)
}

func (mock PricingPlanStoreMock) Delete(entity *database.PricingPlan) error {
	return mock.DeleteFn(entity.ID)
}

func (mock PricingPlanStoreMock) All() ([]database.PricingPlan, error) {
	return mock.AllFn()
}

func (mock PricingPlanStoreMock) ByID(id string) (*database.PricingPlan, error) {
	return mock.ByIDFn(id)
}

func TestNew(t *testing.T) {
	store := PricingPlanStoreMock{}

	testService := service{
		store: &store,
	}

	service := New(&store)

	assert.Equal(t, testService, service)
}

func Test_service_Create(t *testing.T) {
	service := service{
		store: PricingPlanStoreMock{
			CreateFn: func(entity *database.PricingPlan) error {
				return nil
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		assert.NoError(t, service.Create(&PricingPlan{
			AmountUSD: 1,
			Duration:  2,
		}))
	})

	service.store = PricingPlanStoreMock{
		CreateFn: func(entity *database.PricingPlan) error {
			return errors.New("test_error")
		},
	}
	t.Run("Error", func(t *testing.T) {
		assert.Error(t, service.Create(&PricingPlan{}))
	})
}

func Test_service_Update(t *testing.T) {
	service := service{
		store: PricingPlanStoreMock{
			UpdateFn: func(*database.PricingPlan) error {
				return nil
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		assert.NoError(t, service.Update(&PricingPlan{}))
	})

	service.store = PricingPlanStoreMock{
		UpdateFn: func(*database.PricingPlan) error {
			return errors.New("test_error")
		},
	}

	t.Run("Error", func(t *testing.T) {
		assert.Error(t, service.Update(&PricingPlan{}))
	})
}

func Test_service_Delete(t *testing.T) {
	service := service{
		store: PricingPlanStoreMock{
			DeleteFn: func(id string) error {
				return nil
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		assert.NoError(t, service.Delete("id"))
	})

	service.store = PricingPlanStoreMock{
		DeleteFn: func(id string) error {
			return errors.New("test_error")
		},
	}

	t.Run("Error", func(t *testing.T) {
		assert.Error(t, service.Delete("id"))
	})
}

func Test_service_All(t *testing.T) {
	service := service{
		store: PricingPlanStoreMock{
			AllFn: func() ([]database.PricingPlan, error) {
				return []database.PricingPlan{
					{
						ID:        "id",
						AmountUSD: 1,
						Duration:  2,
						MaxUsers:  3,
						UpLimit:   4,
						DownLimit: 5,
						PurgeDays: 6,
					},
					{
						ID:        "id1",
						AmountUSD: 16,
						Duration:  25,
						MaxUsers:  34,
						UpLimit:   43,
						DownLimit: 52,
						PurgeDays: 61,
					},
				}, nil
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		pricingPlans, err := service.All()
		if assert.NoError(t, err) && assert.Equal(t, 2, len(pricingPlans)) {
			assert.Equal(t,
				PricingPlan{
					ID:        "id1",
					AmountUSD: 16,
					Duration:  25,
					MaxUsers:  34,
					UpLimit:   43,
					DownLimit: 52,
					PurgeDays: 61,
				},
				pricingPlans[1],
			)
		}
	})

	service.store = PricingPlanStoreMock{
		AllFn: func() ([]database.PricingPlan, error) {
			return nil, errors.New("test_error")
		},
	}

	t.Run("Error", func(t *testing.T) {
		_, err := service.All()
		assert.Error(t, err)
	})
}
func Test_service_ByID(t *testing.T) {
	service := service{
		store: PricingPlanStoreMock{
			ByIDFn: func(id string) (*database.PricingPlan, error) {
				return &database.PricingPlan{
					ID:        id,
					AmountUSD: 1,
				}, nil
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		pricingPlan, err := service.ByID("id")
		assert.NoError(t, err)
		assert.Equal(t, "id", pricingPlan.ID)
	})

	service.store = PricingPlanStoreMock{
		ByIDFn: func(id string) (*database.PricingPlan, error) {
			return nil, errors.New("test_error")
		},
	}

	t.Run("Error", func(t *testing.T) {
		pricingPlan, err := service.ByID("id")
		assert.Error(t, err)
		assert.Nil(t, pricingPlan)
	})
}
