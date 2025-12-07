package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hekanemre/taxihub/domain"
	"github.com/hekanemre/taxihub/gateway/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}

	return check, msg
}

func Signup(userRepo *helpers.TokenHelper) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user domain.User

		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := validate.Struct(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		count, err := userRepo.UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Println("Error checking email:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error occurred while checking for the email"})
		}

		countPhone, err := userRepo.UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Println("Error checking phone:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error occurred while checking for the phone"})
		}

		if count > 0 || countPhone > 0 {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "this email or phone number already exists"})
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		// Set timestamps and IDs
		now := time.Now().UTC()
		user.Created_at = now
		user.Updated_at = now
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// Generate tokens
		token, refreshToken, err := userRepo.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
		if err != nil {
			log.Println("Error generating tokens:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate tokens"})
		}
		user.Token = &token
		user.Refresh_token = &refreshToken

		// Insert user into MongoDB
		result, insertErr := userRepo.UserCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created: %v", insertErr)
			log.Println(msg)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": msg})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}

func Login(userRepo *helpers.TokenHelper) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user domain.User
		var foundUser domain.User

		// Parse request body
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// Find user by email
		err := userRepo.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "email or password is incorrect"})
		}

		// Verify password
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": msg})
		}

		if foundUser.Email == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}

		// Generate tokens
		token, refreshToken, err := userRepo.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
		if err != nil {
			log.Println("Error generating tokens:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate tokens"})
		}

		// Update tokens in database
		userRepo.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		// Refresh user data after updating tokens
		err = userRepo.UserCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(foundUser)
	}
}
