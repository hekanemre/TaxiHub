package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hekanemre/taxihub/gateway/controllers"
	"github.com/hekanemre/taxihub/gateway/helpers"
)

func AuthRoutes(app *fiber.App, tokenHelper *helpers.TokenHelper) {
	app.Post("/login", controllers.Login(tokenHelper))
	app.Post("/signup", controllers.Signup(tokenHelper))
}
