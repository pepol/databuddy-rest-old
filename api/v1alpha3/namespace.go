package v1alpha3

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// NamespaceService is a service for management of namespaces.
type NamespaceService interface {
	List(prefix string) (NamespaceList, error)
	Has(string) bool
	Get(string) (*NamespaceStatus, error)
	Put(*NamespaceSpec) (*NamespaceStatus, error)
	Delete(string) (*NamespaceStatus, error)
}

// NamespaceList is a list of namespace names returned by list/search operation.
type NamespaceList []string

// NamespaceSpec is the definition of namespace as provided by user.
type NamespaceSpec struct {
	Name        string            `json:",omitempty"`
	Description string            `json:",omitempty"`
	Labels      map[string]string `json:",omitempty"`
}

// NamespaceStatus is the status of namespace as seen by DataBuddy.
type NamespaceStatus struct {
	NamespaceSpec
	CreateIndex string
	UpdateIndex string
	DeleteIndex string `json:",omitempty"`
}

func (ctl *Controller) routeNamespace(router fiber.Router) {
	ns := router.Group("/namespace")

	ns.Get("/", ctl.listNamespaces)
	ns.Get("/+", ctl.getNamespace)
	ns.Put("/+", ctl.putNamespace)
	ns.Delete("/+", ctl.deleteNamespace)
}

// @Summary List all namespaces
// @Tags namespace
// @Accept json
// @Produce json
// @Param prefix query string false "Namespace name prefix" default()
// @Success 200 {object} NamespaceList
// @Router /namespace [get]
// @Description List all namespaces.
func (ctl *Controller) listNamespaces(c *fiber.Ctx) error {
	prefix := c.Query("prefix", "")
	nsList, err := ctl.namespace.List(prefix)
	if err != nil {
		return err
	}

	return c.JSON(nsList)
}

// @Summary Get namespace
// @Tags namespace
// @Accept json
// @Produce json
// @Param name path string true "Namespace name"
// @Success 200 {object} NamespaceStatus
// @Failure 404 {object} RequestError
// @Router /namespace/{name} [get]
// @Description Get namespace by name.
func (ctl *Controller) getNamespace(c *fiber.Ctx) error {
	name := c.Params("+")

	if !ctl.namespace.Has(name) {
		return httpError(c, fiber.StatusNotFound, fmt.Sprintf("namespace '%s' not found", name))
	}

	ns, err := ctl.namespace.Get(name)
	if err != nil {
		return err
	}

	return c.JSON(ns)
}

// @Summary Create/update namespace
// @Tags namespace
// @Accept json
// @Produce json
// @Param name path string true "Namespace name"
// @Param spec body NamespaceSpec true "Namespace fields to update"
// @Success 200 {object} NamespaceStatus
// @Failure 400 {object} RequestError
// @Router /namespace/{name} [put]
// @Description Create the namespace with given name and spec.
// @Description Update fields of given namespace based on body if it already
// @Description exists.
func (ctl *Controller) putNamespace(c *fiber.Ctx) error {
	name := c.Params("+")

	nsSpec := new(NamespaceSpec)
	if err := c.BodyParser(nsSpec); err != nil {
		return httpError(c, fiber.StatusBadRequest, err.Error())
	}

	nsSpec.Name = name

	if ctl.namespace.Has(name) {
		ns, err := ctl.namespace.Get(name)
		if err != nil {
			return err
		}

		if nsSpec.Description == "" {
			nsSpec.Description = ns.Description
		}

		if len(nsSpec.Labels) == 0 {
			nsSpec.Labels = ns.Labels
		}
	}

	nsStored, err := ctl.namespace.Put(nsSpec)
	if err != nil {
		return err
	}

	return c.JSON(nsStored)
}

// @Summary Delete namespace
// @Tags namespace
// @Accept json
// @Produce json
// @Param name path string true "Namespace name"
// @Success 200 {object} NamespaceStatus
// @Failure 404 {object} RequestError
// @Router /namespace/{name} [delete]
// @Description Mark given namespace as deleted.
// @Description All the objects stored within the namespace are scheduled for
// @Description deletion asynchronously. While the namespace is in the process
// @Description of being deleted, GET-ing it will return the object with status
// @Description attribute "DeleteIndex" set to index of the delete operation.
// @Description Once all the contents of the namespace are deleted, GET on
// @Description the namespace will return HTTP 404.
func (ctl *Controller) deleteNamespace(c *fiber.Ctx) error {
	name := c.Params("+")

	if !ctl.namespace.Has(name) {
		return httpError(c, fiber.StatusNotFound, fmt.Sprintf("namespace '%s' not found", name))
	}

	ns, err := ctl.namespace.Delete(name)
	if err != nil {
		return err
	}

	return c.JSON(ns)
}

// RequestError is a struct encompassing HTTP error message into JSON format.
type RequestError struct {
	Error string `json:"error"`
}

func httpError(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(RequestError{
		Error: message,
	})
}
