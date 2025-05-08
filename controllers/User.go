package controllers

import (
	"github.com/Ewan-Reveille/retech/services"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type UserController struct {
	UserService *services.UserService
}

func (uc *UserController) RegisterRoutes(router fiber.Router) {
	router.Post("/register", uc.CreateUser)
	// router.Post("/login", uc.Login)
	// router.Get("/profile", uc.GetProfile)
}

func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	var req services.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	user, err := uc.UserService.Create(&req)
	if err != nil {
		// si conflit clé unique
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "username ou email déjà utilisé",
			})
		}
		// sinon on affiche l’erreur détaillée
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}


	return c.Status(fiber.StatusCreated).JSON(user)
}

func (uc *UserController) Login(c *fiber.Ctx) error {
	var req services.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	user, err := uc.UserService.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	return c.JSON(user)
}

// func (uc *UserController) GetProfile(c *fiber.Ctx) error {
// 	userID := c.Locals("userID").(string)
// 	user, err := uc.UserService.GetProfile(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
// 	}

// 	return c.JSON(user)
// }