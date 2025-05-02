package routes

// import (
// 	"github.com/Ewan-Reveille/retech/controllers"
// 	"github.com/Ewan-Reveille/retech/models"
// 	"github.com/Ewan-Reveille/retech/services"
// 	"github.com/gofiber/fiber/v2"
// 	"gorm.io/gorm"
// )

// func RegisterOrderRoutes(app fiber.Router, db *gorm.DB) {
// 	orderModel := &models.OrderModel{DB: db}
// 	orderService := &services.OrderService{Repo: orderModel}
// 	orderController := &controllers.OrderController{OrderService: orderService}

// 	app.Post("/orders", orderController.CreateOrder)
// 	app.Get("/orders/:id", orderController.GetOrder)
// 	app.Put("/orders/:id", orderController.UpdateOrder)
// 	app.Delete("/orders/:id", orderController.DeleteOrder)
// }
