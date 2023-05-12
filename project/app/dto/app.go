package dto

import (
	"gorm.io/gorm"
)

/*
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
}*/

type Category struct {
	gorm.Model
	Name     string    `json:"name" gorm:"not null"`
	Products []Product `json:"-" gorm:"foreignKey:CategoryID"`
}

type Product struct {
	gorm.Model
	Name         string        `json:"name" gorm:"not null"`
	Description  string        `json:"description" gorm:"not null"`
	Price        float64       `json:"price" gorm:"not null"`
	CategoryID   uint          `json:"categoryID" gorm:"not null"`
	Category     Category      `json:"-"`
	Interactions []Interaction `json:"-" gorm:"foreignKey:ProductID"`
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
	Name     string    `json:"name" gorm:"not null"`
	Email    string    `json:"email" gorm:"not null"`
	Password string    `json:"password" gorm:"not null"`
	Role     string    `json:"role" gorm:"not null"`
	Products []Product `json:"-" gorm:"many2many:interactions;"`
}

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
