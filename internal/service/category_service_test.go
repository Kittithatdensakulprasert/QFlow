package service

import (
	"context"
	"testing"

	"qflow/internal/domain"
	"qflow/internal/repository"

	"github.com/stretchr/testify/assert"
)

func newTestCategoryService() domain.CategoryService {
	return NewCategoryService(&mockCategoryRepository{
		data:   map[uint]domain.Category{},
		nextID: 1,
	})
}

func TestCreateCategory_Success(t *testing.T) {
	service := newTestCategoryService()

	category, err := service.CreateCategory(context.Background(), "ชาบู")

	assert.NoError(t, err)
	assert.NotNil(t, category)
}

func TestCreateCategory_NameRequired(t *testing.T) {
	service := newTestCategoryService()

	category, err := service.CreateCategory(context.Background(), "")

	assert.Nil(t, category)
	assert.Equal(t, ErrCategoryNameRequired, err)
}

func TestCreateCategory_Duplicate(t *testing.T) {
	service := newTestCategoryService()

	_, _ = service.CreateCategory(context.Background(), "ชาบู")
	category, err := service.CreateCategory(context.Background(), "ชาบู")

	assert.Nil(t, category)
	assert.Equal(t, ErrCategoryDuplicate, err)
}

func TestCreateCategory_DuplicateFromRepo(t *testing.T) {
	repo := &mockCategoryRepository{
		createFunc: func(ctx context.Context, category *domain.Category) error {
			return repository.ErrCategoryDuplicate
		},
	}

	service := NewCategoryService(repo)

	category, err := service.CreateCategory(context.Background(), "ชาบู")

	assert.Nil(t, category)
	assert.Equal(t, ErrCategoryDuplicate, err)
}

func TestCreateCategory_ExistsError(t *testing.T) {
	repo := &mockCategoryRepository{
		existsByNameFunc: func(ctx context.Context, name string) (bool, error) {
			return false, errMockRepo
		},
	}

	service := NewCategoryService(repo)

	category, err := service.CreateCategory(context.Background(), "ชาบู")

	assert.Nil(t, category)
	assert.Equal(t, errMockRepo, err)
}

func TestCreateCategory_CreateError(t *testing.T) {
	repo := &mockCategoryRepository{
		createFunc: func(ctx context.Context, category *domain.Category) error {
			return errMockRepo
		},
	}

	service := NewCategoryService(repo)

	category, err := service.CreateCategory(context.Background(), "ชาบู")

	assert.Nil(t, category)
	assert.Equal(t, errMockRepo, err)
}

func TestGetCategories_Success(t *testing.T) {
	service := newTestCategoryService()

	categories, err := service.GetCategories(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, categories)
}

func TestGetCategory_Success(t *testing.T) {
	service := newTestCategoryService()

	created, _ := service.CreateCategory(context.Background(), "ชาบู")

	category, err := service.GetCategory(context.Background(), created.ID)

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, created.ID, category.ID)
	assert.Equal(t, "ชาบู", category.Name)
}

func TestGetCategory_NotFound(t *testing.T) {
	service := newTestCategoryService()

	category, err := service.GetCategory(context.Background(), 999)

	assert.Nil(t, category)
	assert.Equal(t, ErrCategoryNotFound, err)
}

func TestGetCategory_RepoError(t *testing.T) {
	repo := &mockCategoryRepository{
		findByIDFunc: func(ctx context.Context, id uint) (*domain.Category, error) {
			return nil, errMockRepo
		},
	}

	service := NewCategoryService(repo)

	category, err := service.GetCategory(context.Background(), 1)

	assert.Nil(t, category)
	assert.Equal(t, errMockRepo, err)
}

func TestUpdateCategory_Success(t *testing.T) {
	service := newTestCategoryService()

	created, _ := service.CreateCategory(context.Background(), "ชาบู")

	category, err := service.UpdateCategory(context.Background(), created.ID, "ปิ้งย่าง")

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, "ปิ้งย่าง", category.Name)
}

func TestUpdateCategory_NameRequired(t *testing.T) {
	service := newTestCategoryService()

	category, err := service.UpdateCategory(context.Background(), 1, "")

	assert.Nil(t, category)
	assert.Equal(t, ErrCategoryNameRequired, err)
}

func TestUpdateCategory_NotFound(t *testing.T) {
	service := newTestCategoryService()

	category, err := service.UpdateCategory(context.Background(), 999, "ชาบู")

	assert.Nil(t, category)
	assert.Equal(t, ErrCategoryNotFound, err)
}

func TestUpdateCategory_Duplicate(t *testing.T) {
	service := newTestCategoryService()

	c1, _ := service.CreateCategory(context.Background(), "ชาบู")
	_, _ = service.CreateCategory(context.Background(), "ซูชิ")

	res, err := service.UpdateCategory(context.Background(), c1.ID, "ซูชิ")

	assert.Nil(t, res)
	assert.Equal(t, ErrCategoryDuplicate, err)
}

func TestUpdateCategory_DuplicateFromRepo(t *testing.T) {
	repo := &mockCategoryRepository{
		findByIDFunc: func(ctx context.Context, id uint) (*domain.Category, error) {
			return &domain.Category{ID: id, Name: "ชาบู"}, nil
		},
		updateFunc: func(ctx context.Context, category *domain.Category) error {
			return repository.ErrCategoryDuplicate
		},
	}

	service := NewCategoryService(repo)

	category, err := service.UpdateCategory(context.Background(), 1, "ซูชิ")

	assert.Nil(t, category)
	assert.Equal(t, ErrCategoryDuplicate, err)
}

func TestUpdateCategory_RepoError(t *testing.T) {
	repo := &mockCategoryRepository{
		findByIDFunc: func(ctx context.Context, id uint) (*domain.Category, error) {
			return nil, errMockRepo
		},
	}

	service := NewCategoryService(repo)

	category, err := service.UpdateCategory(context.Background(), 1, "ชาบู")

	assert.Nil(t, category)
	assert.Equal(t, errMockRepo, err)
}

func TestUpdateCategory_ExistsError(t *testing.T) {
	repo := &mockCategoryRepository{
		findByIDFunc: func(ctx context.Context, id uint) (*domain.Category, error) {
			return &domain.Category{ID: id, Name: "ชาบู"}, nil
		},
		existsByNameFunc: func(ctx context.Context, name string) (bool, error) {
			return false, errMockRepo
		},
	}

	service := NewCategoryService(repo)

	category, err := service.UpdateCategory(context.Background(), 1, "ซูชิ")

	assert.Nil(t, category)
	assert.Equal(t, errMockRepo, err)
}

func TestUpdateCategory_UpdateError(t *testing.T) {
	repo := &mockCategoryRepository{
		findByIDFunc: func(ctx context.Context, id uint) (*domain.Category, error) {
			return &domain.Category{ID: id, Name: "ชาบู"}, nil
		},
		updateFunc: func(ctx context.Context, category *domain.Category) error {
			return errMockRepo
		},
	}

	service := NewCategoryService(repo)

	category, err := service.UpdateCategory(context.Background(), 1, "ซูชิ")

	assert.Nil(t, category)
	assert.Equal(t, errMockRepo, err)
}

func TestDeleteCategory_Success(t *testing.T) {
	service := newTestCategoryService()

	created, _ := service.CreateCategory(context.Background(), "ชาบู")

	err := service.DeleteCategory(context.Background(), created.ID)

	assert.NoError(t, err)
}

func TestDeleteCategory_NotFound(t *testing.T) {
	service := newTestCategoryService()

	err := service.DeleteCategory(context.Background(), 999)

	assert.Equal(t, ErrCategoryNotFound, err)
}

func TestDeleteCategory_RepoError(t *testing.T) {
	repo := &mockCategoryRepository{
		findByIDFunc: func(ctx context.Context, id uint) (*domain.Category, error) {
			return nil, errMockRepo
		},
	}

	service := NewCategoryService(repo)

	err := service.DeleteCategory(context.Background(), 1)

	assert.Equal(t, errMockRepo, err)
}

func TestDeleteCategory_DeleteError(t *testing.T) {
	repo := &mockCategoryRepository{
		findByIDFunc: func(ctx context.Context, id uint) (*domain.Category, error) {
			return &domain.Category{ID: id, Name: "ชาบู"}, nil
		},
		deleteFunc: func(ctx context.Context, id uint) error {
			return errMockRepo
		},
	}

	service := NewCategoryService(repo)

	err := service.DeleteCategory(context.Background(), 1)

	assert.Equal(t, errMockRepo, err)
}