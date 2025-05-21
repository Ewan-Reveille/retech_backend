package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
	"log"
)

type User struct {
	gorm.Model
	// ID          uuid.UUID `gorm:"type:char(36);primaryKey"` // Tests uniquement
	// ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username    string    `gorm:"uniqueIndex"`
	Email       string    `gorm:"uniqueIndex"`
	Password    string
	PhoneNumber string
	Address     string
	IsVerified  bool
	IsSeller    bool
	IsAdmin     bool
	Rating      float64

	CreatedAt time.Time
	UpdatedAt time.Time

	Products         []Product `gorm:"foreignKey:SellerID"`
	// Orders           []Order   `gorm:"foreignKey:BuyerID"`
	// Reviews          []Review  `gorm:"foreignKey:ReviewerID"`
	// MessagesSent     []Message `gorm:"foreignKey:SenderID"`
	// MessagesReceived []Message `gorm:"foreignKey:ReceiverID"`
	// Reports          []Report  `gorm:"foreignKey:ReporterID"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByPhoneNumber(phoneNumber string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	GetAll() ([]User, error)
}

type UserModel struct {
	DB *gorm.DB
}

func (um *UserModel) Create(user *User) error {
	return um.DB.Create(user).Error
}

func (um *UserModel) GetByID(id uuid.UUID) (*User, error) {
	var user User
	// if err := um.DB.Preload("Products").Preload("Orders").Preload("Reviews").Preload("MessagesSent").Preload("MessagesReceived").Preload("Reports").First(&user, "id = ?", id).Error; err != nil {
	if err := um.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (um *UserModel) GetByUsername(username string) (*User, error) {
	var user User
	// if err := um.DB.Preload("Products").Preload("Orders").Preload("Reviews").Preload("MessagesSent").Preload("MessagesReceived").Preload("Reports").First(&user, "username = ?", username).Error; err != nil {
	if err := um.DB.First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	log.Printf("[UserModel.GetByUsername] trouvé: %+v\n", user)
	return &user, nil
}

func (um *UserModel) GetByEmail(email string) (*User, error) {
	var user User
	// if err := um.DB.Preload("Products").Preload("Orders").Preload("Reviews").Preload("MessagesSent").Preload("MessagesReceived").Preload("Reports").First(&user, "email = ?", email).Error; err != nil {
	if err := um.DB.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (um *UserModel) GetByPhoneNumber(phoneNumber string) (*User, error) {
	var user User
	// if err := um.DB.Preload("Products").Preload("Orders").Preload("Reviews").Preload("MessagesSent").Preload("MessagesReceived").Preload("Reports").First(&user, "phone_number = ?", phoneNumber).Error; err != nil {
	if err := um.DB.First(&user, "phoneNumber = ?", phoneNumber).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (um *UserModel) Update(user *User) error {
	return um.DB.Save(user).Error
}

func (um *UserModel) Delete(id uuid.UUID) error {
	return um.DB.Delete(&User{}, "id = ?", id).Error
}

func (um *UserModel) GetAll() ([]User, error) {
    log.Println("[UserModel.GetAll] début")
    var users []User
    err := um.DB.Find(&users).Error
    if err != nil {
        log.Printf("[UserModel.GetAll] erreur DB.Find(): %v\n", err)
        return nil, err
    }
    log.Printf("[UserModel.GetAll] DB.Find() a renvoyé %d users\n", len(users))
    return users, nil
}

