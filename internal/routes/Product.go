package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Ewan-Reveille/retech/internal/controllers"
	"github.com/Ewan-Reveille/retech/internal/services"
	"github.com/Ewan-Reveille/retech/internal/models"
	"gorm.io/gorm"
)

func RegisterProductRoutes(app fiber.Router, db *gorm.DB) {
	productModel := &models.ProductModel{DB: db}
	productService := &services.ProductService{Repo: productModel}
	productController := &controllers.ProductController{ProductService: productService}

	app.Post("/products", productController.CreateProduct)
	app.Get("/products/:id", productController.GetProduct)
	app.Put("/products/:id", productController.UpdateProduct)
	app.Delete("/products/:id", productController.DeleteProduct)
}
