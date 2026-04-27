package repository

import (
	"context"
	"errors"
	"qflow/internal/domain"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrCategoryRecordNotFound = errors.New("category record not found")

type categoryMemoryRepository struct {
	mu     sync.RWMutex
	data   map[uint]domain.Category
	nextID uint
}

func NewCategoryMemoryRepository() domain.CategoryRepository {
	return &categoryMemoryRepository{
		data:   make(map[uint]domain.Category),
		nextID: 1,
	}
}

func (r *categoryMemoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make([]domain.Category, 0, len(r.data))

	for _, category := range r.data {
		categories = append(categories, category)
	}

	sort.Slice(categories, func(i, j int) bool {
		return categories[i].ID < categories[j].ID
	})

	return categories, nil
}

func (r *categoryMemoryRepository) FindByID(ctx context.Context, id uint) (*domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	category, ok := r.data[id]
	if !ok {
		return nil, ErrCategoryRecordNotFound
	}

	return &category, nil
}

func (r *categoryMemoryRepository) Create(ctx context.Context, category *domain.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	category.ID = r.nextID
	category.CreatedAt = now
	category.UpdatedAt = now

	r.data[category.ID] = *category
	r.nextID++

	return nil
}

func (r *categoryMemoryRepository) Update(ctx context.Context, category *domain.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	oldCategory, ok := r.data[category.ID]
	if !ok {
		return ErrCategoryRecordNotFound
	}

	category.CreatedAt = oldCategory.CreatedAt
	category.UpdatedAt = time.Now()

	r.data[category.ID] = *category

	return nil
}

func (r *categoryMemoryRepository) Delete(ctx context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return ErrCategoryRecordNotFound
	}

	delete(r.data, id)

	return nil
}

func (r *categoryMemoryRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, category := range r.data {
		if strings.EqualFold(category.Name, name) {
			return true, nil
		}
	}

	return false, nil
}