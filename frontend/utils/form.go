package utils

import "github.com/gofiber/fiber/v2"

// PhotoOrder represents the structure of our form data
type PhotoOrder struct {
	FullName  string            `form:"fullName"`
	Location  string            `form:"location"`
	Size      string            `form:"size"`
	PaperType string            `form:"paperType"`
	Email     string            `form:"email"`
	Photos    []*fiber.FormFile `form:"photos"`
}
