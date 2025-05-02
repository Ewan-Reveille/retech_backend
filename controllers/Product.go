package controllers

import (
	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
)

type ProductController struct {
	ProductService *services.ProductService
}

func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {
    var p models.Product
    if err := c.BodyParser(&p); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
    }

    // 1) Création en base
    if err := pc.ProductService.Create(&p); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB failed"})
    }

    // 2) Envoi à la blockchain
    bcProd := services.OnChainProduct{
        ID:          p.ID,
        Title:       p.Title,
        Description: p.Description,
        Price:       p.Price,
        SellerID:    p.SellerID,
	}
	if err := services.SendProduct(bcProd); err != nil {
		log.Printf("blockchain product error: %v", err)
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
