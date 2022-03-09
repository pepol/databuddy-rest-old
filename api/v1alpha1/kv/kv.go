// Package kv contains basic key-value methods.
package kv

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Item contains the store value and key.
type Item struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// Error contains error for any request in KV.
type Error struct {
	Key   string `json:"key"`
	Error string `json:"error"`
}

// GetResponse contains returned value for key.
type GetResponse struct {
	Item
}

// PutResponse contains stored value at key.
type PutResponse struct {
	Item
	Timestamp time.Time `json:"timestamp"`
}

// ItemController is a temporary solution to store 1 KV.
type ItemController struct {
	item *Item
}

// NewItemController returns an ItemController with no KV stored.
func NewItemController() ItemController {
	return ItemController{
		item: nil,
	}
}

// GetItem is a function to retrieve value stored at key
// @Summary Get value stored at key
// @Description Get value stored at key
// @Tags kv
// @Accept json
// @Produce json
// @Param key path string true "Item key"
// @Success 200 {object} GetResponse
// @Failure 404 {object} Error
// @Failure 503 {object} Error
// @Router /item/{key} [get]
// Get value stored for key.
func (ic *ItemController) GetItem(c *fiber.Ctx) error {
	key := c.Params("key")

	if ic.item == nil || ic.item.Key != key {
		return c.Status(http.StatusNotFound).JSON(Error{
			Key:   key,
			Error: "item not found",
		})
	}

	return c.JSON(GetResponse{*ic.item})
}

// PutItem is a function to store value at key
// @Summary Store value for a key
// @Description Store value for a key
// @Tags kv
// @Accept json
// @Accept plain
// @Produce json
// @Param key path string true "Item key"
// @Param value body interface{} true "Value to store"
// @Success 200 {object} PutResponse
// @Failure 400 {object} Error
// @Failure 503 {object} Error
// @Router /item/{key} [post]
// @Router /item/{key} [put]
// Store value at key.
func (ic *ItemController) PutItem(c *fiber.Ctx) error {
	key := c.Params("key")

	item := new(Item)
	item.Key = key

	if c.Is("json") {
		var result map[string]interface{}
		if err := json.Unmarshal(c.Body(), &result); err != nil {
			var resultArr []interface{}
			if err = json.Unmarshal(c.Body(), &resultArr); err != nil {
				return c.Status(http.StatusNotFound).JSON(Error{
					Key:   key,
					Error: err.Error(),
				})
			}
			item.Value = resultArr
		} else {
			item.Value = result
		}
	} else {
		item.Value = string(c.Body())
	}

	ic.item = item
	return c.JSON(PutResponse{*item, time.Now()})
}
