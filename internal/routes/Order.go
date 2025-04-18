package routes

import (
	"github.com/Ewan-Reveille/retech/internal/controllers"
	"github.com/Ewan-Reveille/retech/internal/models"
	"github.com/Ewan-Reveille/retech/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterOrderRoutes(app fiber.Router, db *gorm.DB) {
	orderModel := &models.OrderModel{DB: db}
	orderService := &services.OrderService{Repo: orderModel}
	orderController := &controllers.OrderController{OrderService: orderService}

	app.Post("/orders", orderController.CreateOrder)
	app.Get("/orders/:id", orderController.GetOrder)
	app.Put("/orders/:id", orderController.UpdateOrder)
	app.Delete("/orders/:id", orderController.DeleteOrder)
}