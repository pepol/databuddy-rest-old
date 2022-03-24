// Package v1alpha3 implements the 3rd alpha API version.
package v1alpha3

import "github.com/gofiber/fiber/v2"

// Controller groups all service interfaces for use in Fiber/HTTP routing.
type Controller struct {
	cluster   ClusterService
	namespace NamespaceService
	kv        KeyValueService
	self      NodeService
}

// Config contains the configuration for controller.
type Config struct {
	ClusterService   ClusterService
	KeyValueService  KeyValueService
	NamespaceService NamespaceService
	NodeService      NodeService
}

// NewController initialized controller correctly.
func NewController(config *Config) *Controller {
	ctl := new(Controller)

	ctl.cluster = config.ClusterService
	ctl.kv = config.KeyValueService
	ctl.namespace = config.NamespaceService
	ctl.self = config.NodeService

	return ctl
}

// Route sets up Fiber routing for Controller's services.
func (ctl *Controller) Route(router fiber.Router) {
	ctl.routeKV(router)
	ctl.routeNamespace(router)
	ctl.routeStatus(router)
}
