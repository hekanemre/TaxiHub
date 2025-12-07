package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hekanemre/taxihub/gateway/controllers"
	"github.com/hekanemre/taxihub/gateway/helpers"
	"github.com/hekanemre/taxihub/gateway/middleware"
	"github.com/hekanemre/taxihub/infrastructure"
)

func DriverRoutes(app *fiber.App, driverRepo *infrastructure.MongoRepository, tokenHelper *helpers.TokenHelper) {
	app.Use(middleware.Authenticate(tokenHelper))
	app.Post("/driver/create", controllers.CreateDriver(driverRepo))
	app.Put("/driver/update", controllers.UpdateDriver(driverRepo))
	app.Get("/driver/getall", controllers.GetAllDrivers(driverRepo))
	app.Get("/driver/:id", controllers.GetDriverByID(driverRepo))
	app.Get("driver/getallnearby/:lat/:lon/:taxiType", controllers.GetAllDriversNearby(driverRepo))
}
