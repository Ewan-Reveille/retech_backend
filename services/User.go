package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/Ewan-Reveille/retech/models"
	"github.com/google/uuid"
	"fmt"
	"log"
)

type UserService struct {
	Repo models.UserRepository
}

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (us *UserService) Create(req *CreateUserRequest) (*models.User, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
        return nil, errors.New("tous les champs sont requis")
    }

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User {
		ID: uuid.New(),
		Username: req.Username,
		Email: req.Email,
		Password: string(hash),
	}

	if err := us.Repo.Create(user); err != nil {
		return nil, fmt.Errorf("échec création User: %w", err)
	}
	user.Password = ""
	return user, nil
}

func (us *UserService) Login(req *LoginRequest) (*models.User, error) {
	user, err := us.Repo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("identifiants invalides")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("identifiants invalides")
	}

	user.Password = ""
	return user, nil
}

func (us *UserService) GetAll() ([]models.User, error) {
    log.Println("[UserService.GetAll] début")
    users, err := us.Repo.GetAll()
    if err != nil {
        log.Printf("[UserService.GetAll] erreur Repo.GetAll(): %v\n", err)
        return nil, fmt.Errorf("échec récupération Users: %w", err)
    }
    log.Printf("[UserService.GetAll] Repo.GetAll() a renvoyé %d users\n", len(users))
    return users, nil
}
