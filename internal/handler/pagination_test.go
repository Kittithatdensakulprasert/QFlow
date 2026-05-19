package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestParsePagination_Defaults(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/api/test", nil)

	page, limit, err := parsePagination(c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if page != 1 || limit != 20 {
		t.Fatalf("expected page=1,limit=20 got page=%d limit=%d", page, limit)
	}
}

func TestParsePagination_InvalidPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/api/test?page=0", nil)

	_, _, err := parsePagination(c)
	if err == nil {
		t.Fatal("expected error for invalid page")
	}
}

func TestParsePagination_InvalidLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/api/test?limit=0", nil)

	_, _, err := parsePagination(c)
	if err == nil {
		t.Fatal("expected error for invalid limit")
	}
}

func TestParsePagination_MaxLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/api/test?limit=999", nil)

	_, limit, err := parsePagination(c)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if limit != 100 {
		t.Fatalf("expected limit to be clamped to 100, got %d", limit)
	}
}

func TestPaginateSlice(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	got := paginateSlice(items, 2, 2)
	if len(got) != 2 || got[0] != 3 || got[1] != 4 {
		t.Fatalf("unexpected paginated result: %+v", got)
	}
}

func TestPaginateSlice_OutOfRange(t *testing.T) {
	items := []int{1, 2, 3}
	got := paginateSlice(items, 5, 10)
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %+v", got)
	}
}
