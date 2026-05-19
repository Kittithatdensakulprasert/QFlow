package domain

import "errors"

var (
	ErrProviderCategoryRecordNotFound = errors.New("provider category record not found")
	ErrProviderRecordNotFound         = errors.New("provider record not found")
	ErrProviderZoneRecordNotFound     = errors.New("provider zone record not found")

	ErrQueueZoneRecordNotFound = errors.New("queue zone record not found")
	ErrQueueRecordNotFound     = errors.New("queue record not found")

	ErrCategoryRecordNotFound = errors.New("category record not found")
	ErrCategoryDuplicate      = errors.New("category duplicate")
)
