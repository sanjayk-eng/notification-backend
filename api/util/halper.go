package util

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type User struct {
	CountryCode string `json:"country-code" validate:"required,len=2"`
	Phone       string `json:"phone" validate:"required,len=10"`
	OTP         string `json:"otp" validate:"required,len=6"`
}
type OTPRequest struct {
	CountryCode string `json:"country-code" validate:"required,len=2"`
	Phone       string `json:"phone" validate:"required,len=10"`
}
type AllDetails struct {
	ID           string    `json:"id"`
	Phone        string    `json:"phone"`
	CountryCode  string    `json:"country_code"`
	FullName     string    `json:"fullName"`
	ProfilePhoto string    `json:"profilePhoto,omitempty"`
	Status       string    `json:"status,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type WSMessage struct {
	Phone        string    `json:"phone"`
	FullName     string    `json:"fullName"`
	ProfilePhoto string    `json:"profilePhoto"`
	Message      string    `json:"message"`
	Type         string    `json:"type"`
	CreatedAt    time.Time `json:"created_at"`
}

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()
}

func JsonRespond(c *gin.Context, status int, message any) {
	c.JSON(status, gin.H{"message": message})
}
