package radius

import (
	"strconv"

	"github.com/strusty/worldbit-wifi/radius_database"
	"github.com/strusty/worldbit-wifi/random"
)

type service struct {
	checkStore radius_database.CheckStore
	replyStore radius_database.ReplyStore
}

func New(
	checkStore radius_database.CheckStore,
	replyStore radius_database.ReplyStore,
) Radius {
	return service{
		checkStore: checkStore,
		replyStore: replyStore,
	}
}

func (service service) CreateCredentials(plan PricingPlan) (string, error) {
	usernamePassword := random.String(8)
	if err := service.checkStore.Create(&radius_database.Check{
		Username:  usernamePassword,
		Attribute: "Cleartext-Password",
		Op:        ":=",
		Value:     usernamePassword,
	}); err != nil {
		return "", err
	}

	if err := service.checkStore.Create(&radius_database.Check{
		Username:  usernamePassword,
		Attribute: "Simultaneous-Use",
		Op:        ":=",
		Value:     strconv.FormatInt(plan.MaxUsers, 10),
	}); err != nil {
		return "", err
	}

	if err := service.checkStore.Create(&radius_database.Check{
		Username:  usernamePassword,
		Attribute: "Max-Daily-Session",
		Op:        ":=",
		Value:     strconv.FormatInt(plan.Duration, 10),
	}); err != nil {
		return "", err
	}
	if err := service.replyStore.Create(&radius_database.Reply{
		Username:  usernamePassword,
		Attribute: "WISPr-Bandwidth-Max-Up",
		Op:        "=",
		Value:     strconv.FormatInt(plan.UpLimit, 10),
	}); err != nil {
		return "", err
	}
	if err := service.replyStore.Create(&radius_database.Reply{
		Username:  usernamePassword,
		Attribute: "WISPr-Bandwidth-Max-Down",
		Op:        "=",
		Value:     strconv.FormatInt(plan.DownLimit, 10),
	}); err != nil {
		return "", err
	}

	return usernamePassword, nil
}
