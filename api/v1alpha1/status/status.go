// Package status contains health-check reporting handlers.
package status

import (
	"net/http"
	"os"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

// Response contains status information about the service.
type Response struct {
	Hostname string `json:"hostname,omitempty"`
	Version  string `json:"version,omitempty"`
	Running  bool   `json:"running"`
	Error    string `json:"error,omitempty"`
}

// GetStatus is a function to get service status.
// @Summary Get status of service
// @Description Get status of service, including hostname and version
// @Tags status
// @Produce json
// @Success 200 {object} Response
// @Failure 503 {object} Response
// @Router /status [get]
// Return status of service.
func GetStatus(c *fiber.Ctx) error {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return c.Status(http.StatusServiceUnavailable).JSON(Response{
			Running: false,
			Error:   "unable to get version info",
		})
	}

	hostname, err := os.Hostname()
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(Response{
			Running: false,
			Version: info.Main.Version,
			Error:   err.Error(),
		})
	}

	return c.JSON(Response{
		Running:  true,
		Version:  info.Main.Version,
		Hostname: hostname,
	})
}
