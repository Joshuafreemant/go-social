package helpers

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Joshuafreemant/go-social/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const otpLength = 6

// GenerateOTP generates a numeric OTP
func GenerateOTP() (string, error) {
	bytes := make([]byte, otpLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	otp := ""
	for _, b := range bytes {
		otp += fmt.Sprintf("%d", b%10)
	}
	return otp, nil
}

// SaveOTP stores the OTP in the database with expiration
func SaveOTP(email, otp string, db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	expiresAt := time.Now().Add(10 * time.Minute) // OTP valid for 10 minutes
	otpData := model.OTPData{
		Email:     email,
		OTP:       otp,
		ExpiresAt: expiresAt,
	}

	_, err := db.Collection("otps").UpdateOne(ctx, bson.M{"email": email}, bson.M{
		"$set": otpData,
	}, options.Update().SetUpsert(true))
	return err
}
