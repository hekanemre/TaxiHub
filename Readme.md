# TaxiHub 🚖 (Fiber / MongoDb / Gateway / JWT)

# Overview

A Go-based microservice with MongoDb integration, logging, and API Gateway management

# Prerequisites

Go 1.23.4
Docker and Docker Compose

# Tech Stack

* Language: Go
* Database: MongoDb
* Container: Docker
* Config: Viper
* Logger: Zap
* API Documentation: Swagger
* DesignPattern: Clean Architecture


# Project Structure 
```
├── application
│   ├── driver
│   │   ├── create_driver_handler.go
│   │   ├── driver_service.go
│   │   ├── get_all_driver_handler.go
│   │   ├── get_all_driver_nearby.go
│   │   ├── get_driver_by_plate_handler.go
│   │   ├── get_driver_handler.go
│   │   ├── repository.go
│   │   └── update_driver_handler.go
│   ├── healthcheck
│   │   └── health.go
│   └── error_response.go
├── config
│   ├── config.go
│   └── config.yaml  -- App Configuration File
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── domain
│   ├── driver.go
│   ├── location.go
│   └── user.go
├── gateway
│   ├── controllers
│   │   ├── authController.go
│   │   └── driverController.go
│   ├── helpers
│   │   ├── authHelper.go
│   │   └── tokenHelper.go
│   ├── middleware
│   │   └── authMiddleware.go
│   └── routes
│       ├── authRouter.go
│       └── driverRouter.go
├── infrastructure
│   ├── driverRepository.go
│   └── repository.go
├── log
│   └── log.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
└── Readme.md
```
# Getting Started

1. Clone the repository:
```
git clone https://github.com/hekanemre/TaxiHub
```
2. Start the required services using Docker Compose:
```
docker-compose up -d
```
3. Run the application:
```
go run main.go
```