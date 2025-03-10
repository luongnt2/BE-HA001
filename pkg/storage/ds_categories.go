package storage

import (
	"BE-HA001/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type CategoryStorage interface {
	GetCategoriesByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Category, error)
}

func (s *Storage) GetCategoriesByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Category, error) {
	var results []*model.Category
	var missingIDs []uuid.UUID

	redisKeys := make([]string, len(ids))
	for i, id := range ids {
		redisKeys[i] = s.getCacheKeyCategory(id)
	}

	cachedData, _ := s.cache.MGet(ctx, redisKeys)
	for _, id := range ids {
		key := s.getCacheKeyCategory(id)
		if val, found := cachedData[key]; found {
			var category model.Category
			if err := json.Unmarshal([]byte(val), &category); err == nil {
				results = append(results, &category)
			} else {
				missingIDs = append(missingIDs, id)
			}
		} else {
			missingIDs = append(missingIDs, id)
		}
	}
	if len(missingIDs) == 0 {
		return results, nil
	}

	var categories []*model.Category
	if err := s.DB.WithContext(ctx).Where("id IN ?", missingIDs).Find(&categories).Error; err != nil {
		return nil, err
	}

	cacheData := make(map[string]string)
	for _, category := range categories {
		jsonData, _ := json.Marshal(category)
		cacheData[s.getCacheKeyCategory(category.ID)] = string(jsonData)
		results = append(results, category)
	}

	_ = s.cache.MSet(ctx, cacheData, time.Hour)

	return results, nil
}

func (s *Storage) getCacheKeyCategory(id uuid.UUID) string {
	return fmt.Sprintf("category:%s", id.String())
}
