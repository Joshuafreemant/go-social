package controller

import (
	"context"
	"log"
	"time"

	"github.com/Joshuafreemant/go-social/database"
	"github.com/Joshuafreemant/go-social/helpers"
	"github.com/Joshuafreemant/go-social/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func CreateUsers(c *fiber.Ctx, db *mongo.Database) error {
	users := new(model.Users)
	var ctx = context.Background()

	users.CreatedAt = time.Now()
	users.UpdatedAt = time.Now()

	if err := c.BodyParser(users); err != nil {
		return helpers.ResponseMsg(c, 400, err.Error(), nil)
	}

	if users.Password == "" {
		return helpers.ResponseMsg(c, 400, "Password is required", nil)
	}

	// Check if the email is already taken
	var existingUser model.Users

	err := db.Collection("users").FindOne(ctx, bson.M{"email": users.Email}).Decode(&existingUser)
	if err != nil && err != mongo.ErrNoDocuments {
		return helpers.ResponseMsg(c, 500, "Error checking email", err.Error())
	}
	if existingUser.Email != "" {
		return helpers.ResponseMsg(c, 400, "Email already taken", nil)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users.Password), bcrypt.DefaultCost)
	if err != nil {
		return helpers.ResponseMsg(c, 500, "Failed to hash password", err.Error())
	}

	users.Password = string(hashedPassword)

	if r, err := db.Collection("users").InsertOne(ctx, users); err != nil {
		return helpers.ResponseMsg(c, 500, "Inserted data unsuccessfully", err.Error())
	} else {
		return helpers.ResponseMsg(c, 200, "Inserted data successfully", r)
	}
}

func Login(c *fiber.Ctx, db *mongo.Database) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helpers.ResponseMsg(c, fiber.StatusBadRequest, "Invalid request", err.Error())
	}

	var user model.Users
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return helpers.ResponseMsg(c, fiber.StatusUnauthorized, "Invalid email or password", nil)
	}

	if !helpers.CheckPasswordHash(req.Password, user.Password) {
		return helpers.ResponseMsg(c, fiber.StatusUnauthorized, "Invalid email or password", nil)
	}

	token, err := helpers.GenerateJWT(user.Email)
	if err != nil {
		return helpers.ResponseMsg(c, fiber.StatusInternalServerError, "Failed to generate token", err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
	})

	return helpers.ResponseMsg(c, fiber.StatusOK, "Login successful", fiber.Map{"token": token, "data": user})
}

func ForgotPassword(c *fiber.Ctx) error {
	var ctx = context.Background()
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	type Data struct {
		Email string `json:"email"`
	}

	var req Data
	if err := c.BodyParser(&req); err != nil {
		return helpers.ResponseMsg(c, fiber.StatusBadRequest, "Invalid request", err.Error())
	}

	var user model.Users
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = db.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return helpers.ResponseMsg(c, fiber.StatusUnauthorized, "Invalid email", nil)
	}

	// Generate OTP
	otp, err := helpers.GenerateOTP()
	if err != nil {
		return helpers.ResponseMsg(c, fiber.StatusInternalServerError, "Failed to generate OTP", err.Error())
	}

	// Save OTP to database
	err = helpers.SaveOTP(req.Email, otp, db)
	if err != nil {
		return helpers.ResponseMsg(c, fiber.StatusInternalServerError, "Failed to save OTP", err.Error())
	}

	// Send OTP via email
	err = helpers.SendEmail(req.Email, otp)
	if err != nil {
		return helpers.ResponseMsg(c, fiber.StatusInternalServerError, "Failed to send email", err.Error())
	}

	return helpers.ResponseMsg(c, fiber.StatusOK, "OTP sent successfully", nil)
}

func VerifyOTP(c *fiber.Ctx) error {
	type OTPRequest struct {
		Email string `json:"email"`
		OTP   string `json:"OTP"`
	}
	var req OTPRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.ResponseMsg(c, fiber.StatusBadRequest, "Invalid request", err.Error())
	}

	// Connect to the database
	ctx := context.Background()
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Retrieve OTP data
	var otpData model.OTPData
	err = db.Collection("otps").FindOne(ctx, bson.M{"otp": req.OTP, "email": req.Email}).Decode(&otpData)
	if err != nil {
		return helpers.ResponseMsg(c, fiber.StatusUnauthorized, "Invalid or expired OTP", nil)
	}

	// Compare OTPs
	if otpData.OTP != req.OTP {
		return helpers.ResponseMsg(c, fiber.StatusUnauthorized, "Invalid OTP", nil)
	}

	// Check OTP expiration
	if time.Now().After(otpData.ExpiresAt) {
		return helpers.ResponseMsg(c, fiber.StatusUnauthorized, "OTP has expired", nil)
	}

	// OTP is valid
	return helpers.ResponseMsg(c, fiber.StatusOK, "OTP verified", nil)
}

func ResetPassword(c *fiber.Ctx) error {
	type ResetRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req ResetRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.ResponseMsg(c, fiber.StatusBadRequest, "Invalid request", err.Error())
	}

	if req.Password == "" {
		return helpers.ResponseMsg(c, fiber.StatusBadRequest, "Password is required", nil)
	}

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return helpers.ResponseMsg(c, fiber.StatusInternalServerError, "Failed to hash password", err.Error())
	}

	// Define the update
	update := bson.M{
		"$set": bson.M{
			"password":  string(hashedPassword),
			"updatedAt": time.Now(),
		},
	}

	// Create context
	ctx := context.Background()

	// Update the user's password in the database
	filter := bson.M{"email": req.Email}
	result, err := db.Collection("users").UpdateOne(ctx, filter, update)
	if err != nil {
		return helpers.ResponseMsg(c, fiber.StatusInternalServerError, "Failed to update password", err.Error())
	}

	// Check if the user was found and updated
	if result.MatchedCount == 0 {
		return helpers.ResponseMsg(c, fiber.StatusNotFound, "User not found", nil)
	}

	return helpers.ResponseMsg(c, fiber.StatusOK, "Password updated successfully", result)
}
