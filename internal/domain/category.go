package domain

import (
	"context"
	"time"
)

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]Category, error)
	FindByID(ctx context.Context, id uint) (*Category, error)
	Create(ctx context.Context, category *Category) error
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id uint) error
	ExistsByName(ctx context.Context, name string) (bool, error)
}

type CategoryService interface {
	GetCategories(ctx context.Context) ([]Category, error)
	GetCategory(ctx context.Context, id uint) (*Category, error)
	CreateCategory(ctx context.Context, name string) (*Category, error)
	UpdateCategory(ctx context.Context, id uint, name string) (*Category, error)
	DeleteCategory(ctx context.Context, id uint) error
}
