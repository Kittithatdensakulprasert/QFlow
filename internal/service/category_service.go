package service

import (
	"context"
	"errors"
	"qflow/internal/domain"
	"qflow/internal/repository"
	"strings"
)

var (
	ErrCategoryNameRequired = errors.New("category name is required")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrCategoryDuplicate    = errors.New("category name already exists")
)

type categoryService struct {
	repo domain.CategoryRepository
}

func NewCategoryService(repo domain.CategoryRepository) domain.CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetCategories(ctx context.Context) ([]domain.Category, error) {
	return s.repo.FindAll(ctx)
}

func (s *categoryService) GetCategory(ctx context.Context, id uint) (*domain.Category, error) {
	category, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrCategoryRecordNotFound) {
		return nil, ErrCategoryNotFound
	}

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) CreateCategory(ctx context.Context, name string) (*domain.Category, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, ErrCategoryNameRequired
	}

	exists, err := s.repo.ExistsByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrCategoryDuplicate
	}

	category := &domain.Category{Name: name}

	if err := s.repo.Create(ctx, category); err != nil {
		if errors.Is(err, repository.ErrCategoryDuplicate) {
			return nil, ErrCategoryDuplicate
		}

		return nil, err
	}

	return category, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id uint, name string) (*domain.Category, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, ErrCategoryNameRequired
	}

	category, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrCategoryRecordNotFound) {
		return nil, ErrCategoryNotFound
	}

	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(category.Name, name) {
		exists, err := s.repo.ExistsByName(ctx, name)
		if err != nil {
			return nil, err
		}

		if exists {
			return nil, ErrCategoryDuplicate
		}
	}

	category.Name = name

	if err := s.repo.Update(ctx, category); err != nil {
		if errors.Is(err, repository.ErrCategoryDuplicate) {
			return nil, ErrCategoryDuplicate
		}

		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrCategoryRecordNotFound) {
		return ErrCategoryNotFound
	}

	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
