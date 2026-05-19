package repository

import (
	"context"
	"testing"

	"qflow/internal/domain"
)

func TestCategoryRepo_CreateAndFind(t *testing.T) {
	db := newTestDB(t)
	repo := NewCategoryGormRepository(db)
	ctx := context.Background()

	err := repo.Create(ctx, &domain.Category{Name: "Hospital"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cats, err := repo.FindAll(ctx, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cats) != 1 {
		t.Errorf("expected 1 category, got %d", len(cats))
	}
}

func TestCategoryRepo_Count(t *testing.T) {
	db := newTestDB(t)
	repo := NewCategoryGormRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &domain.Category{Name: "A"})
	repo.Create(ctx, &domain.Category{Name: "B"})

	n, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
}

func TestCategoryRepo_FindByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewCategoryGormRepository(db)
	ctx := context.Background()

	cat := &domain.Category{Name: "Clinic"}
	repo.Create(ctx, cat)

	found, err := repo.FindByID(ctx, cat.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Name != "Clinic" {
		t.Errorf("expected Clinic, got %s", found.Name)
	}
}

func TestCategoryRepo_FindByID_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewCategoryGormRepository(db)
	ctx := context.Background()

	_, err := repo.FindByID(ctx, 9999)
	if err != domain.ErrCategoryRecordNotFound {
		t.Errorf("expected ErrCategoryRecordNotFound, got %v", err)
	}
}

func TestCategoryRepo_Update(t *testing.T) {
	db := newTestDB(t)
	repo := NewCategoryGormRepository(db)
	ctx := context.Background()

	cat := &domain.Category{Name: "Old"}
	repo.Create(ctx, cat)

	cat.Name = "New"
	err := repo.Update(ctx, cat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, _ := repo.FindByID(ctx, cat.ID)
	if found.Name != "New" {
		t.Errorf("expected New, got %s", found.Name)
	}
}

func TestCategoryRepo_Delete(t *testing.T) {
	db := newTestDB(t)
	repo := NewCategoryGormRepository(db)
	ctx := context.Background()

	cat := &domain.Category{Name: "ToDelete"}
	repo.Create(ctx, cat)

	err := repo.Delete(ctx, cat.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = repo.FindByID(ctx, cat.ID)
	if err != domain.ErrCategoryRecordNotFound {
		t.Errorf("expected not found after delete, got %v", err)
	}
}

func TestCategoryRepo_Delete_NotFound(t *testing.T) {
	db := newTestDB(t)
	repo := NewCategoryGormRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, 9999)
	if err != domain.ErrCategoryRecordNotFound {
		t.Errorf("expected ErrCategoryRecordNotFound, got %v", err)
	}
}

func TestCategoryRepo_ExistsByName(t *testing.T) {
	db := newTestDB(t)
	repo := NewCategoryGormRepository(db)
	ctx := context.Background()

	repo.Create(ctx, &domain.Category{Name: "Dental"})

	exists, err := repo.ExistsByName(ctx, "dental") // lowercase
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("expected exists=true")
	}

	exists, err = repo.ExistsByName(ctx, "NonExistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Error("expected exists=false")
	}
}

func TestCategoryRepo_FindAll_Pagination(t *testing.T) {
	db := newTestDB(t)
	repo := NewCategoryGormRepository(db)
	ctx := context.Background()

	for _, name := range []string{"A", "B", "C", "D", "E"} {
		repo.Create(ctx, &domain.Category{Name: name})
	}

	page1, _ := repo.FindAll(ctx, 0, 3)
	if len(page1) != 3 {
		t.Errorf("expected 3, got %d", len(page1))
	}

	page2, _ := repo.FindAll(ctx, 3, 3)
	if len(page2) != 2 {
		t.Errorf("expected 2, got %d", len(page2))
	}
}
