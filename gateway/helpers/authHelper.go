package helpers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// CheckUserType ensures the user has the required role
func CheckUserType(c *fiber.Ctx, role string) error {
	userType, _ := c.Locals("user_type").(string)

	if userType != role {
		return errors.New("unauthorized to access this resource")
	}

	return nil
}

// MatchUserTypeToUid checks both role + ownership
func MatchUserTypeToUid(c *fiber.Ctx, userId string) error {
	userType, _ := c.Locals("user_type").(string)
	uid, _ := c.Locals("uid").(string)

	if userType == "USER" && uid != userId {
		return errors.New("unauthorized to access this resource")
	}

	return CheckUserType(c, userType)
}
