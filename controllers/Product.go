package controllers

import (
	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"

	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type ProductController struct {
	ProductService *services.ProductService
	UserModel      *models.UserModel
}

// CreateProduct creates a new product listing
// @Summary Create a new product
// @Description Create a new product with images and Stripe integration
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param X-User-Username header string true "Seller's username"
// @Param title formData string true "Product title"
// @Param description formData string true "Product description"
// @Param price formData number true "Product price"
// @Param category formData string true "Category ID"
// @Param condition formData string false "Product condition" Enums(new, very good, good, used, fair, unknown)
// @Param images formData file true "Product images"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form"})
	}

	username := c.Get("X-User-Username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "X-User-Username header is required"})
	}

	user, err := pc.UserModel.GetByUsername(username)
	if err != nil {
		log.Printf("error fetching user by username: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user"})
	}

	sellerID := user.ID
	// 1. Parse form fields
	values := form.Value
	title := getFirstValue(values, "title")
	description := getFirstValue(values, "description")
	priceStr := getFirstValue(values, "price")
	categoryIDStr := getFirstValue(values, "category")
	condition := getFirstValue(values, "condition")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		log.Printf("error parsing price: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid price"})
	}

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		log.Printf("error parsing category ID: %v", err)
		log.Printf("%v", categoryID)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	// Default condition if not provided
	if condition == "" {
		condition = string(models.UNKNOWN)
	}

	// Validate required fields
	if title == "" || description == "" || categoryIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title, description, and category, are required"})
	}

	// Parse seller ID
	// sellerID, err := uuid.Parse(sellerIDStr)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid seller_id format"})
	// }

	// Ensure upload directory exists
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		log.Printf("couldn't create upload dir: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server setup error"})
	}

	files := form.File["images"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "At least one image is required"})
	}

	// 3. Build product with ID and images slice
	prodID := uuid.New()
	p := models.Product{
		ID:          prodID,
		Title:       title,
		Description: description,
		Price:       price,
		Condition:   condition,
		CategoryID:  uuid.MustParse(categoryIDStr),
		SellerID:    sellerID,
		Images:      make([]models.ProductImage, 0, len(files)),
	}

	// 4. Save each file and append image records
	for _, fh := range files {
		imgID := uuid.New()
		ext := filepath.Ext(fh.Filename)
		savePath := filepath.Join(uploadDir, fmt.Sprintf("%s%s", imgID.String(), ext))

		if err := c.SaveFile(fh, savePath); err != nil {
			log.Printf("error saving %s: %v", fh.Filename, err)
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": "Failed to save image"})
		}

		p.Images = append(p.Images, models.ProductImage{
			ID:        imgID,
			ProductID: prodID,
			ImageURL:  savePath,
		})
	}

	// 5. Persist to DB and Stripe
	if err := pc.ProductService.Create(&p); err != nil {
		log.Printf("create product error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB or Stripe error"})
	}

	// 6. Fire blockchain event asynchronously
	go func(prod models.Product) {
		bcProd := services.OnChainProduct{
			ID:          prod.ID,
			Title:       prod.Title,
			Description: prod.Description,
			Price:       prod.Price,
			SellerID:    prod.SellerID,
		}
		if err := services.SendProduct(bcProd); err != nil {
			log.Printf("blockchain product error: %v", err)
		}
	}(p)

	return c.Status(fiber.StatusCreated).JSON(p)
}

// GetProduct retrieves a product by ID
// @Summary Get product details
// @Description Get detailed information about a specific product
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [get]
func getFirstValue(values map[string][]string, key string) string {
	if vals := values[key]; len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// GetAllProducts lists all products
// @Summary List all products
// @Description Get a list of all available products
// @Tags Products
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {object} map[string]string
// @Router /products [get]
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


// UpdateProduct updates product details
// @Summary Update a product
// @Description Update an existing product's information
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.Product true "Updated product details"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [put]
func (pc *ProductController) GetAllProducts(c *fiber.Ctx) error {

	products, err := pc.ProductService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve products")
	}

	return c.JSON(products)
}

// DeleteProduct removes a product
// @Summary Delete a product
// @Description Delete a product listing (seller authorization required)
// @Tags Products
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [delete]
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
	log.Printf("Received product: %+v", p)
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
