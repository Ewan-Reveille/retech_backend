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
}

func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form"})
	}
	// 1. Parse form fields
	values := form.Value
	title := getFirstValue(values, "title")
	description := getFirstValue(values, "description")
	priceStr := getFirstValue(values, "price")
	categoryIDStr := getFirstValue(values, "category")
	sellerIDStr := getFirstValue(values, "seller_id")
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
	if title == "" || description == "" || categoryIDStr == "" || sellerIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title, description, category, and seller_id are required"})
	}

	// Parse seller ID
	sellerID, err := uuid.Parse(sellerIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid seller_id format"})
	}

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

func getFirstValue(values map[string][]string, key string) string {
	if vals := values[key]; len(vals) > 0 {
		return vals[0]
	}
	return ""
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

func (pc *ProductController) GetAllProducts(c *fiber.Ctx) error {

	products, err := pc.ProductService.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve products")
	}

	return c.JSON(products)
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
