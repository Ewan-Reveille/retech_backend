package controllers

import (
	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MessageController struct {
	MessageService *services.MessageService
}


// POST /messages
// @Summary Create a new message
// @Description Create a new message between users
// @Tags Messages
// @Accept json
// @Produce json
// @Param message body models.Message true "Message object"
// @Success 201 {object} models.Message
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /messages [post]
func (mc *MessageController) CreateMessage(c *fiber.Ctx) error {
	var message models.Message
	if err := c.BodyParser(&message); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := mc.MessageService.Create(&message); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create message",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(message)
}

// GET /messages/:id
// @Summary Get a message by ID
// @Description Retrieve a specific message by its ID
// @Tags Messages
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} models.Message
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /messages/{id} [get]
func (mc *MessageController) GetMessage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
	}

	message, err := mc.MessageService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Not found")
	}
	return c.JSON(message)
}

// PUT /messages/:id
// @Summary Update a message
// @Description Update an existing message's content
// @Tags Messages
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Param message body models.Message true "Updated message object"
// @Success 200 {object} models.Message
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /messages/{id} [put]
func (mc *MessageController) UpdateMessage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
	}

	var message models.Message
	if err := c.BodyParser(&message); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	message.ID = id

	if err := mc.MessageService.Update(&message); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update message",
		})
	}

	return c.JSON(message)
}

// DELETE /messages/:id
// @Summary Delete a message
// @Description Delete a message by ID
// @Tags Messages
// @Param id path string true "Message ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /messages/{id} [delete]
func (mc *MessageController) DeleteMessage(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
	}

	if err := mc.MessageService.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete message",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}