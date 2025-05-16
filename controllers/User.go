package controllers

import (
	"github.com/Ewan-Reveille/retech/services"
	"github.com/gofiber/fiber/v2"
	"strings"
	"github.com/google/uuid"
	"log"
)

type UserController struct {
	UserService *services.UserService
}

func (uc *UserController) RegisterRoutes(router fiber.Router) {
	router.Post("/register", uc.CreateUser)
	router.Post("/login", uc.Login)
	router.Get("/user", uc.GetAllUsers)
}

func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	var req services.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid body"})
	}

	user, err := uc.UserService.Create(&req)
	if err != nil {

		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "username ou email déjà utilisé",
			})
		}

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

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	IsSeller bool      `json:"isSeller"`
}


func (uc *UserController) GetAllUsers(c *fiber.Ctx) error {
    log.Println("[GetAllUsers] début")
    users, err := uc.UserService.GetAll()
    if err != nil {
        log.Printf("[GetAllUsers] erreur service.GetAll(): %v\n", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch users"})
    }
    log.Printf("[GetAllUsers] service.GetAll() a renvoyé %d users\n", len(users))

    var response []UserResponse
    for _, user := range users {
        response = append(response, UserResponse{
            ID:       user.ID,
            Username: user.Username,
            Email:    user.Email,
            IsSeller: user.IsSeller,
        })
    }
    log.Printf("[GetAllUsers] renvoi réponse JSON: %+v\n", response)
    return c.JSON(response)
}
