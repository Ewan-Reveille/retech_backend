package controllers

import (
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
    "github.com/Ewan-Reveille/retech/internal/models"
    "github.com/Ewan-Reveille/retech/internal/services"
)

type ProductController struct {
	ProductService *services.ProductService
}

// POST /products
func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {
	var p models.Product
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := pc.ProductService.Create(&p); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(p)
}


func (pc *ProductController) GetProduct(c *fiber.Ctx) error {
    idStr := c.Params("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
    }

    product, err := pc.ProductService.GetByID(id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).SendString("Not found")
    }

    return c.JSON(product)
}


// PUT /products/:id
func (pc *ProductController) UpdateProduct(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
	}

	var p models.Product
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if _, err := pc.ProductService.GetByID(id); err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Product not found")
	}

	if err := pc.ProductService.Update(&p); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update product")
	}

	return c.JSON(p)
}

// DELETE /products/:id
func (pc *ProductController) DeleteProduct(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
	}

	product, err := pc.ProductService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Product not found")
	}

	if err := pc.ProductService.Delete(id, product.SellerID); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete product")
	}

	return c.JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}
