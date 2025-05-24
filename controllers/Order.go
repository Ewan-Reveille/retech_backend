package controllers

import (
	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderController struct {
	OrderService *services.OrderService
}

// POST /orders
// @Summary Create a new order
// @Description Create a new order with Stripe payment integration
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body models.Order true "Order details"
// @Success 201 {object} map[string]interface{} "Returns created order and Stripe client secret"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func (oc *OrderController) CreateOrder(c *fiber.Ctx) error {
	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	orderResult, clientSecret, err := oc.OrderService.CreateWithStripe(order.BuyerID, order.ProductID, order.ProductID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}

	// Tu peux aussi renvoyer le `clientSecret` dans ta réponse :
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"order":         orderResult,
		"client_secret": clientSecret,
	})
}

// GET /orders/:id
// @Summary Get order by ID
// @Description Retrieve order details by order ID
// @Tags Orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [get]
func (oc *OrderController) GetOrder(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
	}

	order, err := oc.OrderService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Not found")
	}

	return c.JSON(order)
}

// PUT /orders/:id
// @Summary Update order
// @Description Update existing order details
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param order body models.Order true "Updated order details"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [put]
func (oc *OrderController) UpdateOrder(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
	}

	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	order.ID = id

	if err := oc.OrderService.Update(&order); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	return c.JSON(order)
}

// DELETE /orders/:id
// @Summary Delete order
// @Description Delete an order by ID
// @Tags Orders
// @Param id path string true "Order ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{id} [delete]
func (oc *OrderController) DeleteOrder(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
	}

	if err := oc.OrderService.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete order",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

