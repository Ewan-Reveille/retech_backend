package controllers

import (
	"log"

	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"fmt"
	"os"
	"path/filepath"
	"strconv"

	// "github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82"
	"bytes"
    "io"
    "mime/multipart"
    "net/http"
)

type ProductController struct {
	ProductService *services.ProductService
	UserModel      *models.UserModel
}

// CreateProduct creates a new product listing
// @Summary Create a new product
// @Description Create a new product with images and Stripe integration, and enrich images via API tagging
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
	// 1. Récupération du formulaire
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse form"})
	}

	// 2. Vérification et récupération de l'utilisateur via l'en-tête X-User-Username
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

	// 3. Lecture des champs texte du formulaire
	values := form.Value
	title := getFirstValue(values, "title")
	description := getFirstValue(values, "description")
	priceStr := getFirstValue(values, "price")
	categoryIDStr := getFirstValue(values, "category")
	condition := getFirstValue(values, "condition")

	// Validation prix
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		log.Printf("error parsing price: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid price"})
	}

	// Validation catégorie
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		log.Printf("error parsing category ID: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	// Condition par défaut si non fournie
	if condition == "" {
		condition = string(models.UNKNOWN)
	}
	// Vérification des champs obligatoires
	if title == "" || description == "" || categoryIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title, description, and category are required"})
	}

	// 4. Création du dossier d’uploads si nécessaire
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		log.Printf("couldn't create upload dir: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server setup error"})
	}

	// 5. Récupération des fichiers uploadés
	files := form.File["images"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "At least one image is required"})
	}

	// 6. Construction de la structure Product (sans images pour l'instant)
	prodID := uuid.New()
	p := models.Product{
		ID:          prodID,
		Title:       title,
		Description: description,
		Price:       price,
		Condition:   condition,
		CategoryID:  categoryID,
		SellerID:    sellerID,
		Images:      make([]models.ProductImage, 0, len(files)),
	}

	// 7. Pour chaque fichier, on :
	//    a) le sauvegarde en local
	//    b) l’envoie à l’API de tagging pour récupérer la légende
	//    c) crée une entrée ProductImage avec l’Alt rempli
	for _, fh := range files {
		imgID := uuid.New()
		ext := filepath.Ext(fh.Filename)
		savePath := filepath.Join(uploadDir, fmt.Sprintf("%s%s", imgID.String(), ext))

		// a) Sauvegarde locale
		if err := c.SaveFile(fh, savePath); err != nil {
			log.Printf("error saving %s: %v", fh.Filename, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save image"})
		}

		// b) Envoi du fichier à l’API de tagging
		altText, tagErr := fetchImageCaption(savePath)
		if tagErr != nil {
			// Si l’API échoue, on log et on continue sans alt ou avec alt vide
			log.Printf("error fetching caption for %s: %v", savePath, tagErr)
			altText = ""
		}

		// c) Ajout au slice d’images du produit
		p.Images = append(p.Images, models.ProductImage{
			ID:        imgID,
			ProductID: prodID,
			ImageURL:  savePath,
			Alt:       altText,
		})
	}

	// 8. Persistance en base et intégration Stripe
	if err := pc.ProductService.Create(&p); err != nil {
		log.Printf("create product error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB or Stripe error"})
	}

	// 9. Envoi asynchrone vers la blockchain
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

func (pc *ProductController) CreatePaymentIntent(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	product, err := pc.ProductService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(product.Price * 100)), // in cents
		Currency: stripe.String("eur"),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create PaymentIntent"})
	}

	return c.JSON(fiber.Map{
		"clientSecret": pi.ClientSecret,
	})
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

// fetchImageCaption envoie le fichier image à l’API de tagging et retourne le header X-image-caption
func fetchImageCaption(filePath string) (string, error) {
	// 1. Ouverture du fichier
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	// 2. Création d’un buffer multipart/form-data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 3. Création du champ "image" dans le formulaire
	part, err := writer.CreateFormFile("image", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("cannot create form file: %w", err)
	}

	// 4. Copie du contenu du fichier dans le champ
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("cannot copy file data: %w", err)
	}

	// 5. Fermeture du writer pour terminer le multipart
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("cannot close writer: %w", err)
	}

	// 6. Construction de la requête HTTP POST
	aiaToken := os.Getenv("AIA_TOKEN")
	url := fmt.Sprintf("https://tagging.aia-handicap.com/analyze-image/%s", aiaToken)
	req, err := http.NewRequest(http.MethodPost, url, &requestBody)
	if err != nil {
		return "", fmt.Errorf("cannot create HTTP request: %w", err)
	}
	// Ajout de l’en-tête Content-Type pour indiquer le multipart
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 7. Envoi de la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// 8. Lecture du header X-image-caption
	caption := resp.Header.Get("X-image-caption")
	// On peut valider le code HTTP (optionnel)
	if resp.StatusCode != http.StatusOK {
		return caption, fmt.Errorf("non-200 status: %d", resp.StatusCode)
	}

	return caption, nil
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
