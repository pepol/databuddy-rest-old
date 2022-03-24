package v1alpha3

import "github.com/gofiber/fiber/v2"

// RequestError is a struct encompassing HTTP error message into JSON format.
type RequestError struct {
	Error string `json:"error"`
}

func httpError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(RequestError{
		Error: message,
	})
}
