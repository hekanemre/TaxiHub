package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hekanemre/taxihub/application/healthcheck"
	"github.com/hekanemre/taxihub/config"
	"github.com/hekanemre/taxihub/gateway/helpers"
	"github.com/hekanemre/taxihub/gateway/routes"
	"github.com/hekanemre/taxihub/infrastructure"
	"github.com/hekanemre/taxihub/log"
	"go.uber.org/zap"
)

type Request any
type Response any

type HandlerInterface[Req Request, Res Response] interface {
	Handle(ctx context.Context, req *Req) (*Res, error)
}

func handle[Req Request, Res Response](h HandlerInterface[Req, Res]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req Req
		if err := c.BodyParser(&req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
			zap.L().Error("Failed to parse request body", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if err := c.ParamsParser(&req); err != nil {
			zap.L().Error("Failed to parse request params", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request params",
			})
		}

		if err := c.QueryParser(&req); err != nil {
			zap.L().Error("Failed to parse request query", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request query",
			})
		}

		if err := c.ReqHeaderParser(&req); err != nil {
			zap.L().Error("Failed to parse request headers", zap.Error(err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request headers",
			})
		}

		//the reason why we use c.UserContext() is to propagate the context from fiber to our handler
		//so if the request is cancelled or times out, our handler can also be aware of it
		//this is important for long running requests or when we have a timeout set in fiber
		//so our handler can stop processing if the client has disconnected
		ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
		defer cancel()

		res, err := h.Handle(ctx, &req)
		if err != nil {
			zap.L().Error("Handler error", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error: " + err.Error(),
			})
		}

		return c.JSON(res)
	}
}

func main() {
	appConfig := config.Read()
	log.Init()
	defer zap.L().Sync()

	zap.L().Info("Starting server...")

	app := fiber.New(fiber.Config{
		//timeout is a must for production, fiber cuts long requests
		IdleTimeout:  appConfig.IdleTimeout,
		ReadTimeout:  appConfig.ReadTimeout,
		WriteTimeout: appConfig.WriteTimeout,
		Concurrency:  256 * 1024,
	})

	app.Use(func(c *fiber.Ctx) error {
		// log request details
		zap.L().Info("Request", zap.String("method", c.Method()), zap.String("path", c.Path()))
		return c.Next()
	})

	userRepo, err := infrastructure.NewMongoRepository("users")
	if err != nil {
		zap.L().Error("Failed to connect to MongoDB (users)", zap.Error(err))
		os.Exit(1)
	}

	driverRepo, err := infrastructure.NewMongoRepository("drivers")
	if err != nil {
		zap.L().Error("Failed to connect to MongoDB (drivers)", zap.Error(err))
		os.Exit(1)
	}
	tokenHelper := helpers.NewTokenHelper(userRepo)

	healthCheckHandler := healthcheck.NewHealthCheckHandler()
	app.Get("/health", handle[healthcheck.HealthCheckRequest, healthcheck.HealthCheckResponse](healthCheckHandler))

	routes.AuthRoutes(app, tokenHelper)
	routes.DriverRoutes(app, driverRepo, tokenHelper)

	zap.L().Info("Server started on port", zap.String("port", appConfig.Port))

	if err := app.Listen(fmt.Sprintf(":%s", appConfig.Port)); err != nil {
		zap.L().Error("Failed to start server", zap.Error(err))
		os.Exit(1)
	}
}
