package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"qflow/internal/domain"

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
			return domain.ErrCategoryDuplicate
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

	categories, total, err := service.GetCategories(context.Background(), 1, 20)

	assert.NoError(t, err)
	assert.NotNil(t, categories)
	assert.Equal(t, int64(len(categories)), total)
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
			return domain.ErrCategoryDuplicate
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

var errMockRepo = errors.New("mock repository error")

type mockCategoryRepository struct {
	data   map[uint]domain.Category
	nextID uint

	findAllFunc      func(ctx context.Context, offset, limit int) ([]domain.Category, error)
	countFunc        func(ctx context.Context) (int64, error)
	findByIDFunc     func(ctx context.Context, id uint) (*domain.Category, error)
	createFunc       func(ctx context.Context, category *domain.Category) error
	updateFunc       func(ctx context.Context, category *domain.Category) error
	deleteFunc       func(ctx context.Context, id uint) error
	existsByNameFunc func(ctx context.Context, name string) (bool, error)
}

func (m *mockCategoryRepository) FindAll(ctx context.Context, offset, limit int) ([]domain.Category, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx, offset, limit)
	}

	categories := make([]domain.Category, 0, len(m.data))
	for _, category := range m.data {
		categories = append(categories, category)
	}

	return categories, nil
}

func (m *mockCategoryRepository) Count(ctx context.Context) (int64, error) {
	if m.countFunc != nil {
		return m.countFunc(ctx)
	}
	return int64(len(m.data)), nil
}

func (m *mockCategoryRepository) FindByID(ctx context.Context, id uint) (*domain.Category, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}

	category, ok := m.data[id]
	if !ok {
		return nil, domain.ErrCategoryRecordNotFound
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
