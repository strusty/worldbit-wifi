package routing

import (
	"net/http"
	"strings"

	"git.sfxdx.ru/crystalline/wi-fi-backend/jwt"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/admins"
	"git.sfxdx.ru/crystalline/wi-fi-backend/services/pricing_plans"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

var jsonValidator *validator.Validate

type AdminRouter struct {
	adminService        admins.Admins
	pricingPlansService pricing_plans.PricingPlans
}

func NewAdminRouter(adminService admins.Admins, pricingPlansService pricing_plans.PricingPlans) AdminRouter {
	return AdminRouter{
		adminService:        adminService,
		pricingPlansService: pricingPlansService,
	}
}

func (router AdminRouter) Register(group *echo.Group) {
	group.POST("/login", router.login)
	group.POST("/plans", router.createPlan)
	group.PUT("/plans/:id", router.updatePlan)
	group.DELETE("/plans/:id", router.deletePlan)
	group.PUT("/password", router.changePassword)
}

func (router AdminRouter) login(context echo.Context) error {
	request := new(admins.LoginRequest)
	if err := context.Bind(request); err != nil {
		return err
	}

	response, err := router.adminService.Login(*request)
	if err != nil {
		return err
	}

	return context.JSON(http.StatusOK, response)
}

func (router AdminRouter) createPlan(context echo.Context) error {
	request := new(pricing_plans.PricingPlanValidator)
	if err := context.Bind(request); err != nil {
		return err
	}

	jsonValidator = validator.New()
	if err := jsonValidator.Struct(request); err != nil {
		return err
	}

	pricingPlanEntity, err := request.GetPricingPlan()
	if err != nil {
		return err
	}

	if err := router.pricingPlansService.Create(pricingPlanEntity); err != nil {
		return err
	}

	return context.JSON(http.StatusOK, pricingPlanEntity)
}

func (router AdminRouter) updatePlan(context echo.Context) error {
	request := new(pricing_plans.PricingPlanValidator)
	if err := context.Bind(request); err != nil {
		return err
	}

	jsonValidator = validator.New()
	if err := jsonValidator.Struct(request); err != nil {
		return err
	}

	if context.Param("id") != request.ID {
		return errors.New("Path id doesn't match request body id.")
	}

	pricingPlanEntity, err := request.GetPricingPlan()
	if err != nil {
		return err
	}

	if err := router.pricingPlansService.Update(pricingPlanEntity); err != nil {
		return err
	}

	return context.JSON(http.StatusOK, pricingPlanEntity)
}

func (router AdminRouter) deletePlan(context echo.Context) error {
	if err := router.pricingPlansService.Delete(context.Param("id")); err != nil {
		return err
	}

	return context.NoContent(http.StatusOK)
}

func (router AdminRouter) changePassword(context echo.Context) error {
	request := new(admins.ChangePasswordRequest)
	if err := context.Bind(request); err != nil {
		return err
	}

	if strings.TrimSpace(request.NewPassword) == "" {
		return errors.New("Password cannot be an empty string")
	}

	userID, err := jwt.GetUserIDFromJWT(context)
	if err != nil {
		return err
	}

	return router.adminService.ChangePassword(userID, *request)
}
