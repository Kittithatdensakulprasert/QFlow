package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPage  = 1
	defaultLimit = 20
	maxLimit     = 100
)

func parsePagination(c *gin.Context) (int, int, error) {
	page := defaultPage
	limit := defaultLimit

	if rawPage := c.Query("page"); rawPage != "" {
		parsedPage, err := strconv.Atoi(rawPage)
		if err != nil || parsedPage <= 0 {
			return 0, 0, fmt.Errorf("page must be a positive integer")
		}
		page = parsedPage
	}

	if rawLimit := c.Query("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err != nil || parsedLimit <= 0 {
			return 0, 0, fmt.Errorf("limit must be a positive integer")
		}
		if parsedLimit > maxLimit {
			parsedLimit = maxLimit
		}
		limit = parsedLimit
	}

	return page, limit, nil
}

func paginateSlice[T any](items []T, page, limit int) []T {
	start := (page - 1) * limit
	if start >= len(items) {
		return []T{}
	}

	end := start + limit
	if end > len(items) {
		end = len(items)
	}

	return items[start:end]
}
