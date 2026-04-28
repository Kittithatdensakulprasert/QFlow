package service

import (
	"context"
	"errors"
	"qflow/internal/domain"
	"strings"
)

var errMockRepo = errors.New("mock repository error")

type mockCategoryRepository struct {
	data   map[uint]domain.Category
	nextID uint

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

	categories := make([]domain.Category, 0, len(m.data))
	for _, category := range m.data {
		categories = append(categories, category)
	}

	return categories, nil
}

func (m *mockCategoryRepository) FindByID(ctx context.Context, id uint) (*domain.Category, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}

	category, ok := m.data[id]
	if !ok {
		return nil, errors.New("category not found")
	}

	return &category, nil
}

func (m *mockCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, category)
	}

	if m.data == nil {
		m.data = map[uint]domain.Category{}
	}

	if m.nextID == 0 {
		m.nextID = 1
	}

	category.ID = m.nextID
	m.data[category.ID] = *category
	m.nextID++

	return nil
}

func (m *mockCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, category)
	}

	if m.data == nil {
		m.data = map[uint]domain.Category{}
	}

	m.data[category.ID] = *category
	return nil
}

func (m *mockCategoryRepository) Delete(ctx context.Context, id uint) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}

	delete(m.data, id)
	return nil
}

func (m *mockCategoryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	if m.existsByNameFunc != nil {
		return m.existsByNameFunc(ctx, name)
	}

	for _, category := range m.data {
		if strings.EqualFold(category.Name, name) {
			return true, nil
		}
	}

	return false, nil
}