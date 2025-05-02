package controllers

import (
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/services"
)

type PaymentController struct {
	PaymentService *services.PaymentService
}

// POST /payments
func (pc *PaymentController) CreatePayment(c *fiber.Ctx) error {
	var p models.Payment
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := pc.PaymentService.Create(&p); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create payment",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(p)
}
// GET /payments/:id
func (pc *PaymentController) GetPayment(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
	}

	payment, err := pc.PaymentService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Not found")
	}

	return c.JSON(payment)
}

// PUT /payments/:id
// func (pc *PaymentController) UpdatePayment(c *fiber.Ctx) error {
// 	idStr := c.Params("id")
// 	id, err := uuid.Parse(idStr)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
// 	}

// 	var p models.Payment
// 	if err := c.BodyParser(&p); err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
// 	}

// 	p.ID = id

// 	if err := pc.PaymentService.Update(&p); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to update payment",
// 		})
// 	}

// 	return c.JSON(p)
// }

// // DELETE /payments/:id
// func (pc *PaymentController) DeletePayment(c *fiber.Ctx) error {
// 	idStr := c.Params("id")
// 	id, err := uuid.Parse(idStr)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
// 	}

// 	if err := pc.PaymentService.Delete(id); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to delete payment",
// 		})
// 	}

// 	return c.SendStatus(fiber.StatusNoContent)
// }

// // GET /payments/buyer/:buyerId
// func (pc *PaymentController) GetPaymentsByBuyerID(c *fiber.Ctx) error {
// 	buyerIdStr := c.Params("buyerId")
// 	buyerId, err := uuid.Parse(buyerIdStr)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
// 	}

// 	payments, err := pc.PaymentService.GetByBuyerID(buyerId)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).SendString("Not found")
// 	}

// 	return c.JSON(payments)
// }

// // GET /payments/seller/:sellerId
// func (pc *PaymentController) GetPaymentsBySellerID(c *fiber.Ctx) error {
// 	sellerIdStr := c.Params("sellerId")
// 	sellerId, err := uuid.Parse(sellerIdStr)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
// 	}

// 	payments, err := pc.PaymentService.GetBySellerID(sellerId)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).SendString("Not found")
// 	}

// 	return c.JSON(payments)
// }

// // GET /payments/order/:orderId
// func (pc *PaymentController) GetPaymentsByOrderID(c *fiber.Ctx) error {
// 	orderIdStr := c.Params("orderId")
// 	orderId, err := uuid.Parse(orderIdStr)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
// 	}

// 	payments, err := pc.PaymentService.GetByOrderID(orderId)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).SendString("Not found")
// 	}

// 	return c.JSON(payments)
// }

// // GET /payments/product/:productId
// func (pc *PaymentController) GetPaymentsByProductID(c *fiber.Ctx) error {
// 	productIdStr := c.Params("productId")
// 	productId, err := uuid.Parse(productIdStr)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
// 	}

// 	payments, err := pc.PaymentService.GetByProductID(productId)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).SendString("Not found")
// 	}

// 	return c.JSON(payments)
// }

// // GET /payments/shipping/:shippingId
// func (pc *PaymentController) GetPaymentsByShippingID(c *fiber.Ctx) error {
// 	shippingIdStr := c.Params("shippingId")
// 	shippingId, err := uuid.Parse(shippingIdStr)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
// 	}

// 	payments, err := pc.PaymentService.GetByShippingID(shippingId)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).SendString("Not found")
// 	}

// 	return c.JSON(payments)
// }

