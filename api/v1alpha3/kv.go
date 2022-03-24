package v1alpha3

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	flagsBitSize = 8
	ttlBitSize   = 64
)

// KeyValueService is a service for management of key-value store.
type KeyValueService interface {
	List(namespace, prefix string) (KeyList, error)
	Has(namespace, key string) bool
	Get(namespace, key string) (*KVItem, error)
	Put(namespace, key string, value []byte, flags uint8, ttl uint64) (*KVItem, error)
	Delete(namespace, key string) (*KVItem, error)
}

// KeyList is a list of stored keys returned by list/search operation.
type KeyList []string

// KVItem is the information about an item stored by KV store.
type KVItem struct {
	Key         string
	CreateIndex string
	UpdateIndex string
	Flags       byte
	ExpiresAt   uint64
	// Value is encoded using base64.
	Value []byte
}

func (ctl *Controller) routeKV(router fiber.Router) {
	ns := router.Group("/kv")

	ns.Get("/", ctl.listKeys)
	ns.Get("/+", ctl.getKey)
	ns.Put("/+", ctl.putKey)
	ns.Delete("/+", ctl.deleteKey)
}

// @Summary List all keys
// @Tags kv
// @Accept json
// @Produce json
// @Param prefix query string false "Key prefix" default()
// @Param namespace query string false "Namespace" default(default)
// @Success 200 {object} KeyList
// @Router /kv [get]
// @Description List all keys.
func (ctl *Controller) listKeys(c *fiber.Ctx) error {
	prefix := c.Query("prefix", "")
	ns := c.Query("namespace", "default")

	keyList, err := ctl.kv.List(ns, prefix)
	if err != nil {
		return err
	}

	return c.JSON(keyList)
}

// @Summary Get key
// @Tags namespace
// @Accept json
// @Produce json
// @Param key path string true "Key"
// @Param namespace query string false "Namespace" default(default)
// @Param raw query bool false "Return only value" default(false)
// @Success 200 {object} KVItem
// @Failure 400 {object} RequestError "Returned when 'raw' parameter is not parseable as boolean"
// @Failure 404 {object} RequestError "Returned when either key or namespace doesn't exist"
// @Router /kv/{key} [get]
// @Description Get key value.
func (ctl *Controller) getKey(c *fiber.Ctx) error {
	key := c.Params("+")
	ns := c.Query("namespace", "default")

	raw, err := strconv.ParseBool(c.Query("raw", "false"))
	if err != nil {
		return httpError(c, fiber.StatusBadRequest, err.Error())
	}

	if !ctl.namespace.Has(ns) {
		return httpError(c, fiber.StatusNotFound, fmt.Sprintf("namespace '%s' not found", ns))
	}

	if !ctl.kv.Has(ns, key) {
		return httpError(c, fiber.StatusNotFound, fmt.Sprintf("key '%s' not found in namespace '%s'", key, ns))
	}

	kvItem, err := ctl.kv.Get(ns, key)
	if err != nil {
		return err
	}

	if raw {
		return c.Send(kvItem.Value)
	}
	return c.JSON(kvItem)
}

// @Summary Put key
// @Tags kv
// @Accept plain
// @Accept octet-stream
// @Produce json
// @Param key path string true "Key"
// @Param namespace query string false "Namespace" default(default)
// @Param flags query byte false "User-defined metadata" default(0)
// @Param ttl query uint64 false "Time-To-Live (in seconds), 0 means the item won't expire" default(0)
// @Param value body string true "Value to store"
// @Success 200 {object} bool
// @Failure 400 {object} RequestError "Returned when no value is provided"
// @Router /kv/{key} [put]
// @Description Store the provided value under key.
// @Description If namespace doesn't exist, it gets created.
func (ctl *Controller) putKey(c *fiber.Ctx) error {
	key := c.Params("+")

	ns := c.Query("namespace", "default")
	if !ctl.namespace.Has(ns) {
		nsSpec := new(NamespaceSpec)
		nsSpec.Name = ns
		nsSpec.Description = ns
		_, err := ctl.namespace.Put(nsSpec)
		if err != nil {
			return err
		}
	}

	flags64, err := strconv.ParseUint(c.Query("flags", "0"), 0, flagsBitSize)
	if err != nil {
		return httpError(c, fiber.StatusBadRequest, err.Error())
	}
	flags := uint8(flags64)

	ttl, err := strconv.ParseUint(c.Query("ttl", "0"), 0, ttlBitSize)
	if err != nil {
		return httpError(c, fiber.StatusBadRequest, err.Error())
	}

	value := c.Body()
	if len(value) == 0 {
		return httpError(c, fiber.StatusBadRequest, "no value to store")
	}

	kv, err := ctl.kv.Put(ns, key, value, flags, ttl)
	if err != nil {
		return err
	}

	return c.JSON(kv)
}

// @Summary Delete key
// @Tags kv
// @Accept json
// @Produce json
// @Param key path string true "Key"
// @Param namespace query string false "Namespace" default(default)
// @Success 200 {object} KVItem
// @Failure 404 {object} RequestError "Returned when either key or namespace doesn't exist"
// @Router /kv/{key} [delete]
// @Description Delete provided key.
func (ctl *Controller) deleteKey(c *fiber.Ctx) error {
	key := c.Params("+")

	ns := c.Query("namespace", "default")
	if !ctl.namespace.Has(ns) {
		return httpError(c, fiber.StatusNotFound, fmt.Sprintf("namespace '%s' not found", ns))
	}

	if !ctl.kv.Has(ns, key) {
		return httpError(c, fiber.StatusNotFound, fmt.Sprintf("key '%s' not found in namespace '%s'", key, ns))
	}

	kv, err := ctl.kv.Delete(ns, key)
	if err != nil {
		return err
	}

	return c.JSON(kv)
}
