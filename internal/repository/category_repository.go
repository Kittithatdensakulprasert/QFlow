package repository

import (
	"context"
	"errors"
	"qflow/internal/domain"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type categoryGormRepository struct {
	db *gorm.DB
}

func NewCategoryGormRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryGormRepository{db: db}
}

func (r *categoryGormRepository) FindAll(ctx context.Context, offset, limit int) ([]domain.Category, error) {
	var categories []domain.Category

	err := r.db.WithContext(ctx).
		Order("id ASC").
		Offset(offset).
		Limit(limit).
		Find(&categories).Error

	return categories, err
}

func (r *categoryGormRepository) Count(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&domain.Category{}).Count(&total).Error
	return total, err
}

func (r *categoryGormRepository) FindByID(ctx context.Context, id uint) (*domain.Category, error) {
	var category domain.Category

	err := r.db.WithContext(ctx).First(&category, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrCategoryRecordNotFound
	}

	return &category, err
}

func (r *categoryGormRepository) Create(ctx context.Context, category *domain.Category) error {
	err := r.db.WithContext(ctx).Create(category).Error
	if isUniqueViolation(err) {
		return domain.ErrCategoryDuplicate
	}

	return err
}

func (r *categoryGormRepository) Update(ctx context.Context, category *domain.Category) error {
	err := r.db.WithContext(ctx).Save(category).Error
	if isUniqueViolation(err) {
		return domain.ErrCategoryDuplicate
	}

	return err
}

func (r *categoryGormRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.Category{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrCategoryRecordNotFound
	}

	return nil
}

func (r *categoryGormRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&domain.Category{}).
		Where("LOWER(name) = ?", strings.ToLower(name)).
		Count(&count).Error

	return count > 0, err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
