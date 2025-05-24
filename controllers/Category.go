package controllers

import (
	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CategoryController struct {
	CategoryService *services.CategoryService
}

// POST /categories
// @Summary Create a new category
// @Description Create a new product category
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body models.SwaggerCategory true "Category object"
// @Success 201 {object} models.SwaggerCategory
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories [post]
func (cc *CategoryController) CreateCategory(c *fiber.Ctx) error {
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}

	if err := cc.CategoryService.Create(&category); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create category",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// GET /categories/:id
// @Summary Get a category by ID
// @Description Get detailed information about a specific category
// @Tags Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id} [get]
func (cc *CategoryController) GetCategory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
	}

	category, err := cc.CategoryService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Not found")
	}

	return c.JSON(category)
}

// PUT /categories/:id
// @Summary Update a category
// @Description Update an existing category's information
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body models.Category true "Updated category object"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id} [put]
func (cc *CategoryController) UpdateCategory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
	}

	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	category.ID = id

	if err := cc.CategoryService.Update(&category); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update category",
		})
	}

	return c.JSON(category)
}

// DELETE /categories/:id
// @Summary Delete a category
// @Description Delete a category by ID
// @Tags Categories
// @Param id path string true "Category ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id} [delete]
func (cc *CategoryController) DeleteCategory(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
	}

	if err := cc.CategoryService.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete category",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GET /categories
// @Summary Get all categories
// @Description Get a list of all product categories
// @Tags Categories
// @Produce json
// @Success 200 {array} models.Category
// @Failure 500 {object} map[string]string
// @Router /categories [get]
func (cc *CategoryController) GetAllCategories(c *fiber.Ctx) error {
	categories, err := cc.CategoryService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve categories",
		})
	}

	return c.JSON(categories)
}

// GET /categories/:id/products
// @Summary Get products by category
// @Description Get all products belonging to a specific category
// @Tags Categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {array} models.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id}/products [get]
func (cc *CategoryController) GetCategoryProducts(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid UUID")
	}

	category, err := cc.CategoryService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Not found")
	}

	return c.JSON(category.Products)
}
