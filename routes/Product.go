package routes

import (
	"github.com/Ewan-Reveille/retech/controllers"
	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterProductRoutes(app fiber.Router, db *gorm.DB, stripeClient services.StripeClient) {
	productModel := &models.ProductModel{DB: db}
	productService := &services.ProductService{Repo: productModel, DB: db, StripeClient: stripeClient}
	userModel := &models.UserModel{DB: db}
	productController := &controllers.ProductController{
        ProductService: productService,
        UserModel:      userModel,
    }

	app.Post("/products", productController.CreateProduct)
	app.Get("/products/:id", productController.GetProduct)
	app.Get("/products", productController.GetAllProducts)
	app.Put("/products/:id", productController.UpdateProduct)
	app.Delete("/products/:id", productController.DeleteProduct)
}
