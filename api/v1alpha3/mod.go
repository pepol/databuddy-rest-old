// Package v1alpha3 implements the 3rd alpha API version.
package v1alpha3

import "github.com/gofiber/fiber/v2"

// Controller groups all service interfaces for use in Fiber/HTTP routing.
type Controller struct {
	namespace NamespaceService
}

// NewController initialized controller correctly.
func NewController(nsSvc NamespaceService) *Controller {
	ctl := new(Controller)

	ctl.namespace = nsSvc

	return ctl
}

// Route sets up Fiber routing for Controller's services.
func (ctl *Controller) Route(router fiber.Router) {
	ctl.routeNamespace(router)
}
