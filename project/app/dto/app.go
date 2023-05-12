package dto

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Role     string `json:"role" gorm:"not null"`
}
type Category struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
}

type Product struct {
	gorm.Model
	Name         string  `json:"name" gorm:"not null"`
	Description  string  `json:"description" gorm:"not null"`
	Price        float64 `json:"price" gorm:"not null"`
	CategoryID   uint    `json:"categoryID" gorm:"not null"`
	Interactions uint    `json:"interactions" `
}

type Interaction struct {
	gorm.Model
	UserID     uint `json:"userId"`
	ProductID  uint `json:"productId"`
	CategoryID uint `json:"categoryId"`
}

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
