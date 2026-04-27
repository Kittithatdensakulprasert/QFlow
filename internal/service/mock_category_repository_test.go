package service

import (
	"context"
	"errors"
	"qflow/internal/domain"
)

var errMockRepo = errors.New("mock repository error")

type mockCategoryRepository struct {
	findAllFunc      func(ctx context.Context) ([]domain.Category, error)
	findByIDFunc     func(ctx context.Context, id uint) (*domain.Category, error)
	createFunc       func(ctx context.Context, category *domain.Category) error
	updateFunc       func(ctx context.Context, category *domain.Category) error
	deleteFunc       func(ctx context.Context, id uint) error
	existsByNameFunc func(ctx context.Context, name string) (bool, error)
}

func (m *mockCategoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx)
	}
	return []domain.Category{}, nil
}

func (m *mockCategoryRepository) FindByID(ctx context.Context, id uint) (*domain.Category, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return &domain.Category{ID: id, Name: "ชาบู"}, nil
}

func (m *mockCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, category)
	}
	return nil
}

func (m *mockCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, category)
	}
	return nil
}

func (m *mockCategoryRepository) Delete(ctx context.Context, id uint) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockCategoryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	if m.existsByNameFunc != nil {
		return m.existsByNameFunc(ctx, name)
	}
	return false, nil
}