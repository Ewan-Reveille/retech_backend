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

// CreateUser creates a new user
// @Summary Register a new user
// @Description Create a new user account
// @Tags Users
// @Accept json
// @Produce json
// @Param user body services.CreateUserRequest true "User registration data"
// @Success 201 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
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

// Login authenticates a user
// @Summary User login
// @Description Authenticate user and return user details
// @Tags Users
// @Accept json
// @Produce json
// @Param credentials body services.LoginRequest true "Login credentials"
// @Success 200 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
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


// GetAllUsers retrieves all users
// @Summary List all users
// @Description Get a list of all registered users
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} UserResponse
// @Failure 500 {object} map[string]string
// @Router /user [get]
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
