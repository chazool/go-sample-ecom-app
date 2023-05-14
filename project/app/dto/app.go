package dto

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name     string    `json:"name" gorm:"not null"`
	Products []Product `json:"-" gorm:"foreignKey:CategoryID"`
}

type Product struct {
	gorm.Model
	Name          string   `json:"name" gorm:"not null"`
	Description   string   `json:"description" gorm:"not null"`
	Price         float64  `json:"price" gorm:"not null"`
	CategoryID    uint     `json:"categoryID" gorm:"not null"`
	Category      Category `json:"-"`
	Interactions  uint     `json:"interactions"`
	WeightedScore float64  `json:"-" gorm:"-"`
}

type Interaction struct {
	gorm.Model
	UserID     uint     `json:"userId"`
	User       User     `json:"user"`
	ProductID  uint     `json:"productId"`
	Product    Product  `json:"-"`
	CategoryID uint     `json:"categoryId"`
	Category   Category `json:"-"`
}

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Role     Role   `json:"role" gorm:"not null"`
	Products []Product `json:"-" gorm:"many2many:interactions;"`
}

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type UserRole struct {
	ID   uint
	Role Role
}

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
	RoleSystem Role = "system"
)
