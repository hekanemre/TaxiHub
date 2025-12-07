package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	application "github.com/hekanemre/taxihub/application/driver"
	"github.com/hekanemre/taxihub/infrastructure"
)

func CreateDriver(driverRepo *infrastructure.MongoRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		createDriverHandler := application.NewCreateDriverHandler(driverRepo)
		GetDriverByPlateHandler := application.NewGetDriverByPlateHandler(driverRepo)

		var validationReq application.GetDriverByPlateRequest
		if err := c.BodyParser(&validationReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		validationRes, validationErr := GetDriverByPlateHandler.Handle(c.UserContext(), &validationReq)

		if validationErr == nil && validationRes != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Driver with the same plate already exists"})
		}

		var req application.CreateDriverRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		res, err := createDriverHandler.Handle(c.UserContext(), &req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(res)
	}
}

func UpdateDriver(driverRepo *infrastructure.MongoRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		updateDriverHandler := application.NewUpdateDriverHandler(driverRepo)
		var req application.UpdateDriverRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}
		res, err := updateDriverHandler.Handle(c.UserContext(), &req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func GetAllDrivers(driverRepo *infrastructure.MongoRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req application.GetAllFilterRequest

		page := c.QueryInt("page", 1)
		pageSize := c.QueryInt("pageSize", 20)

		req.Page = page
		req.PageSize = pageSize

		getAllDriversHandler := application.NewGetAllDriverHandler(driverRepo)

		res, err := getAllDriversHandler.Handle(c.UserContext(), &req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func GetAllDriversNearby(driverRepo *infrastructure.MongoRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		getAllDriversNearbyHandler := application.NewGetAllDriverNearbyHandler(driverRepo)

		latStr := c.Params("lat")
		if latStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing 'lat' query parameter"})
		}
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid 'lat' query parameter"})
		}

		lonStr := c.Params("lon")
		if lonStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing 'lon' query parameter"})
		}
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid 'lon' query parameter"})
		}

		taxiType := c.Params("taxiType")
		if taxiType == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing 'taxiType' query parameter"})
		}

		req := &application.GetAllDriverNearbyRequest{
			Lat:      lat,
			Lon:      lon,
			TaxiType: taxiType,
		}

		res, err := getAllDriversNearbyHandler.Handle(c.UserContext(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func GetDriverByID(driverRepo *infrastructure.MongoRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		getDriverByIDHandler := application.NewGetDriverHandler(driverRepo)

		var req application.GetDriverRequest

		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing id parameter"})
		}
		req.ID = id

		res, err := getDriverByIDHandler.Handle(c.UserContext(), &req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}
