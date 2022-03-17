package v1alpha2

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get collection by name
// @Tags collection
// @Accept json
// @Produce json
// @Param namespace path string true "Name of the collection namespace"
// @Param collection path string true "Name of the collection to retrieve"
// @Success 200 {object} Collection
// @Failure 404 {object} RequestError
// @Router /{namespace}/{collection} [get]
// @Description Retrieve detailed information about collection by name.
func (ctl *Controller) getColl(c *fiber.Ctx) error {
	ns := c.Params("namespace")
	name := c.Params("collection")

	coll, err := ctl.CollSvc.Get(ns, name)
	switch err.(type) {
	case nil:
		return c.JSON(coll)
	case *ErrorNotFound:
		return c.Status(fiber.StatusNotFound).JSON(RequestError{
			Error: fmt.Sprintf("collection '%s/%s' not found", ns, name),
		})
	default:
		return err
	}
}

// @Summary Set collection
// @Tags collection
// @Accept json
// @Produce json
// @Param namespace path string true "Name of the collection namespace"
// @Param collection path string true "Name of the collection"
// @Param spec body Collection true "The collection object"
// @Success 200 {object} Collection
// @Failure 400 {object} RequestError
// @Failure 404 {object} RequestError
// @Router /{namespace}/{collection} [put]
// @Description Modify namespace with "collection" name (path parameter) to match
// @Description the provided collection object. Create the collection if it does not exist.
// @Description The name provided in path and name in request body (if set) MUST
// @Description be the same.
// @Description When creating a new collection, namespace MUST exist, otherwise
// @Description "Not Found" (404) is returned.
func (ctl *Controller) setColl(c *fiber.Ctx) error {
	collection := new(Collection)

	if err := c.BodyParser(collection); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(RequestError{
			Error: err.Error(),
		})
	}

	if name := c.Params("collection"); name != collection.Name {
		return c.Status(fiber.StatusBadRequest).JSON(RequestError{
			Error: fmt.Sprintf("collection names must be the same '%s' (path) != '%s' (body)", name, collection.Name),
		})
	}

	if ns := c.Params("namespace"); ns != collection.Namespace {
		return c.Status(fiber.StatusBadRequest).JSON(RequestError{
			Error: fmt.Sprintf("namespace names must be the same '%s' (path) != '%s' (body)", ns, collection.Namespace),
		})
	}

	coll, err := ctl.CollSvc.Set(collection)
	switch err.(type) {
	case nil:
		return c.JSON(coll)
	case *ErrorNotFound:
		return c.Status(fiber.StatusNotFound).JSON(RequestError{
			Error: fmt.Sprintf("namespace '%s' does not exist", collection.Namespace),
		})
	default:
		return err
	}
}

// @Summary Delete collection
// @Tags collection
// @Accept json
// @Produce json
// @Param namespace path string true "Name of the collection namespace"
// @Param collection path string true "Name of the collection"
// @Success 200 {object} Collection
// @Failure 404 {object} RequestError
// @Router /{namespace}/{collection} [delete]
// @Description Delete provided collection.
func (ctl *Controller) deleteColl(c *fiber.Ctx) error {
	ns := c.Params("namespace")
	name := c.Params("collection")

	coll, err := ctl.CollSvc.Delete(ns, name)
	switch err.(type) {
	case nil:
		return c.JSON(coll)
	case *ErrorNotFound:
		return c.Status(fiber.StatusNotFound).JSON(RequestError{
			Error: fmt.Sprintf("collection '%s/%s' does not exist", ns, name),
		})
	default:
		return err
	}
}

func (ctl *Controller) deleteAllCollections(ns string, colls CollectionList) error {
	for _, c := range colls {
		_, err := ctl.CollSvc.Delete(ns, c.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
