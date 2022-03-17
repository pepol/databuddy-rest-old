package v1alpha2

import (
	"github.com/gofiber/fiber/v2"
)

// Controller is the API handler for API operations.
type Controller struct {
	NsSvc   NamespaceService
	CollSvc CollectionService
}

// Route sets up routing on given Fiber instance.
func (ctl *Controller) Route(router fiber.Router) {
	router.Get("/", ctl.listNS)

	router.Get("/:namespace", ctl.getNS)
	router.Put("/:namespace", ctl.setNS)
	router.Delete("/:namespace", ctl.deleteNS)

	router.Get("/:namespace/:collection", ctl.getColl)
	router.Put("/:namespace/:collection", ctl.setColl)
	router.Delete("/:namespace/:collection", ctl.deleteColl)
}
