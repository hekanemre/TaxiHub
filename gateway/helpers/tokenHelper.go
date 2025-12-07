package helpers

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hekanemre/taxihub/infrastructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Uid        string `json:"uid"`
	User_type  string `json:"user_type"`
	jwt.RegisteredClaims
}

var SECRET_KEY = "your_secret_key_here"

type TokenHelper struct {
	UserCollection *mongo.Collection
}

func NewTokenHelper(repo *infrastructure.MongoRepository) *TokenHelper {
	return &TokenHelper{
		UserCollection: repo.DB.Collection(repo.Collection),
	}
}

func (t *TokenHelper) GenerateAllTokens(
	email, firstName, lastName, userType, uid string,
) (string, string, error) {

	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		User_type:  userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Println("JWT token generation error:", err)
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Println("JWT refresh token generation error:", err)
		return "", "", err
	}

	return token, refreshToken, nil
}

func (t *TokenHelper) ValidateToken(signedToken string) (*SignedDetails, string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		return nil, err.Error()
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, "the token is invalid"
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, "the token is expired"
	}

	return claims, ""
}

func (t *TokenHelper) UpdateAllTokens(signedToken, signedRefreshToken, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateObj := primitive.D{
		{Key: "token", Value: signedToken},
		{Key: "refresh_token", Value: signedRefreshToken},
		{Key: "updated_at", Value: time.Now()},
	}

	filter := bson.M{"user_id": userId}
	upsert := true
	opts := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := t.UserCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opts)
	if err != nil {
		log.Println("MongoDB token update error:", err)
	}
}
