package routes

import (
	"github.com/Ewan-Reveille/retech/controllers"
	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterCategoryRoutes(app fiber.Router, db *gorm.DB) {
	categoryModel := &models.CategoryModel{DB: db}
	categoryService := &services.CategoryService{Repo: categoryModel}
	categoryController := &controllers.CategoryController{CategoryService: categoryService}

	app.Post("/categories", categoryController.CreateCategory)
	app.Get("/categories/:id", categoryController.GetCategory)
	app.Put("/categories/:id", categoryController.UpdateCategory)
	app.Delete("/categories/:id", categoryController.DeleteCategory)
	app.Get("/categories", categoryController.GetAllCategories)
	app.Get("/categories/:id/products", categoryController.GetCategoryProducts)
}
