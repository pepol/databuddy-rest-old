package v1alpha2

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// NamespaceController is the API handler for namespace operations.
type NamespaceController struct {
	Service NamespaceService
}

// Route sets up routing on given Fiber instance.
func (nc *NamespaceController) Route(router fiber.Router) {
	router.Get("/", nc.listNS)

	router.Get("/:namespace", nc.getNS)
	router.Put("/:namespace", nc.setNS)
	router.Delete("/:namespace", nc.deleteNS)
}

// @Summary List accessible namespaces
// @Tags namespace
// @Accept json
// @Produce json
// @Param prefix query string false "Prefix for namespace names" default("")
// @Success 200 {object} NamespaceList
// @Router / [get]
// @Description Retrieve a list of all namespaces in DataBuddy system.
// @Description If RBAC is enabled, the list returned contains only namespaces
// @Description visible to the authenticated user.
// @Description Optional query parameter "prefix" can be provided to return
// @Description only namespaces with the given prefix.
func (nc *NamespaceController) listNS(c *fiber.Ctx) error {
	prefix := c.Query("prefix", "")
	namespaces, err := nc.Service.List(prefix)
	if err != nil {
		return err
	}
	return c.JSON(namespaces)
}

// @Summary Get namespace by name
// @Tags namespace
// @Accept json
// @Produce json
// @Param namespace path string true "Name of the namespace to retrieve"
// @Success 200 {object} Namespace
// @Failure 400 {object} RequestError
// @Failure 404 {object} RequestError
// @Router /{namespace} [get]
// @Description Retrieve detailed information about namespace by name.
func (nc *NamespaceController) getNS(c *fiber.Ctx) error {
	name := c.Params("namespace")

	ns, err := nc.Service.Get(name)
	switch err.(type) {
	case nil:
		return c.JSON(ns)
	case *ErrorNotFound:
		return c.Status(fiber.StatusNotFound).JSON(RequestError{
			Error: fmt.Sprintf("namespace '%s' not found", name),
		})
	default:
		return err
	}
}

// @Summary Set namespace
// @Tags namespace
// @Accept json
// @Produce json
// @Param namespace path string true "Name of the namespace"
// @Param spec body Namespace true "The namespace object"
// @Success 200 {object} Namespace
// @Failure 400 {object} RequestError
// @Router /{namespace} [put]
// @Description Modify namespace with "name" (path parameter) to match
// @Description the provided namespace object. Create namespace if does not exist.
// @Description The name provided in path and name in request body (if set) MUST
// @Description be the same.
func (nc *NamespaceController) setNS(c *fiber.Ctx) error {
	namespace := new(Namespace)

	if err := c.BodyParser(namespace); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(RequestError{
			Error: err.Error(),
		})
	}

	if name := c.Params("namespace"); name != namespace.Name {
		return c.Status(fiber.StatusBadRequest).JSON(RequestError{
			Error: fmt.Sprintf("namespace names must be the same '%s' (path) != '%s' (body)", name, namespace.Name),
		})
	}

	ns, err := nc.Service.Set(namespace)
	if err != nil {
		return err
	}

	return c.JSON(ns)
}

// @Summare Delete namespace
// @Tags namespace
// @Accept json
// @Produce json
// @Param namespace path string true "Name of the namespace"
// @Success 200 {object} Namespace
// @Failure 404 {object} RequestError
// @Router /{namespace} [delete]
// @Description Delete provided namespace.
// @Description This method also deletes all collections that are part of the namespace.
func (nc *NamespaceController) deleteNS(c *fiber.Ctx) error {
	name := c.Params("name")

	ns, err := nc.Service.Delete(name)
	switch err.(type) {
	case nil:
		return c.JSON(ns)
	case *ErrorNotFound:
		return c.Status(fiber.StatusNotFound).JSON(RequestError{
			Error: fmt.Sprintf("namespace '%s' does not exist", name),
		})
	default:
		return err
	}
}
