package dto

import (
	"time"

	"gorm.io/gorm"
)

// define a product model
//
//	type Product struct {
//		gorm.Model
//		Name        string  `json:"name"`
//		Description string  `json:"description"`
//		Price       float64 `json:"price"`
//		Category    string  `json:"category"`
//	}
// type Product struct {
// 	ID              uint      `json:"id" gorm:"primaryKey"`
// 	Name            string    `json:"name"`
// 	Description     string    `json:"description"`
// 	Price           float64   `json:"price"`
// 	CreatedAt       time.Time `json:"created_at"`
// 	UpdatedAt       time.Time `json:"updated_at"`
// 	RelatedProducts []Product `json:"related_products,omitempty" gorm:"-"`
// }

// define a user model
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

// define a web tracking model
type Tracking struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
	Action    string `json:"action"`
}

// define a recommendation model
type Recommendation struct {
	ProductID     uint    `json:"product_id"`
	WeightedScore float64 `json:"weighted_score"`
}

// type Interaction struct {
// 	ID         uint   `gorm:"primaryKey"`
// 	UserID     uint   `gorm:"not null"`
// 	ProductID  uint   `gorm:"not null"`
// 	ActionType string `gorm:"not null"`
// 	CreatedAt  time.Time
// }

type Interaction struct {
	gorm.Model
	UserID      uint
	ProductID   uint
	Interaction string
}

type Product struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	Price       float64   `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func GetDst() []interface{} {
	return []interface{}{&Product{}, &User{}, &Tracking{}, &Recommendation{}, &Interaction{}}

}
