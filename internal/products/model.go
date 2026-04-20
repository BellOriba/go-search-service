package products

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID           uuid.UUID      `json:"id" validate:"required"`
	SKU          string         `json:"sku" validate:"required,alphanum,max=20"`
	Name         string         `json:"name" validate:"required,min=3,max=100"`
	Slug         string         `json:"slug" validate:"required,lowercase"`
	Description  string         `json:"description" validate:"max=1000"`
	Price        int64          `json:"price" validate:"required,gt=0"`
	Stock        int            `json:"stock" validate:"min=0"`
	CategoryID   uuid.UUID      `json:"category_id" validate:"required"`
	CategoryName string         `json:"category_name"`
	IsFeatured   bool           `json:"is_featured"`
	Images       []ProductImage `json:"images" validate:"dive"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type ProductImage struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	Path      string    `json:"path" validate:"required"`
	Original  string    `json:"original_url"`
	Thumbnail string    `json:"thumbnail_url"`
	IsPrimary bool      `json:"is_primary"`
}

type CreateImageRequest struct {
	Path      string `json:"path"`
	Original  string `json:"original_url"`
	Thumbnail string `json:"thumbnail_url"`
	IsPrimary bool   `json:"is_primary"`
}

type CreateProductRequest struct {
	SKU         string               `json:"sku" validate:"required,alphanum,max=20"`
	Name        string               `json:"name" validate:"required,min=3,max=100"`
	Slug        string               `json:"slug" validate:"required,lowercase"`
	Description string               `json:"description" validate:"max=1000"`
	Price       int64                `json:"price" validate:"required,gt=0"`
	Stock       int                  `json:"stock" validate:"min=0"`
	CategoryID  uuid.UUID            `json:"category_id" validate:"required"`
	IsFeatured  bool                 `json:"is_featured"`
	Images      []CreateImageRequest `json:"images"`
}

type UpdateProductImageRequest struct {
	Path      string `json:"path" validate:"required"`
	Original  string `json:"original_url" validate:"required,url"`
	Thumbnail string `json:"thumbnail_url" validate:"required,url"`
}

type ProductIndex struct {
	ID          string `json:"id"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Category    string `json:"category"`
	Image       string `json:"image"`
	IsFeatured  bool   `json:"is_featured"`
	CreatedAt   int64  `json:"created_at"`
}

type User struct {
	ID uuid.UUID `json:"id"`
	Email string `json:"email"`
	PasswordHash string `json:"-"`
	Role string `json:"role"`
}

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

