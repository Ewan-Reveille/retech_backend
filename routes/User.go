package routes

import (
	"github.com/Ewan-Reveille/retech/controllers"
	"github.com/Ewan-Reveille/retech/models"
	"github.com/Ewan-Reveille/retech/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterUserRoutes(app fiber.Router, db *gorm.DB) {
	userModel := &models.UserModel{DB: db}
	userService := &services.UserService{Repo: userModel}
	userController := &controllers.UserController{UserService: userService}

	app.Post("/register", userController.CreateUser)
}

func LoginUserRoutes(app fiber.Router, db *gorm.DB) {
	userModel := &models.UserModel{DB: db}
	userService := &services.UserService{Repo: userModel}
	userController := &controllers.UserController{UserService: userService}

	app.Post("/login", userController.Login)
}