package storage

import (
	"BE-HA001/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type SupplierStorage interface {
	GetSupplierByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Supplier, error)
}

func (s *Storage) GetSupplierByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Supplier, error) {
	var results []*model.Supplier
	var missingIDs []uuid.UUID

	redisKeys := make([]string, len(ids))
	for i, id := range ids {
		redisKeys[i] = s.getCacheKeySupplier(id)
	}

	cachedData, _ := s.cache.MGet(ctx, redisKeys)
	for _, id := range ids {
		key := s.getCacheKeySupplier(id)
		if val, found := cachedData[key]; found {
			var supplier model.Supplier
			if err := json.Unmarshal([]byte(val), &supplier); err == nil {
				results = append(results, &supplier)
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

	var suppliers []*model.Supplier
	if err := s.DB.WithContext(ctx).Where("id IN ?", missingIDs).Find(&suppliers).Error; err != nil {
		return nil, err
	}

	cacheData := make(map[string]string)
	for _, supplier := range suppliers {
		jsonData, _ := json.Marshal(supplier)
		cacheData[s.getCacheKeySupplier(supplier.ID)] = string(jsonData)
		results = append(results, supplier)
	}

	_ = s.cache.MSet(ctx, cacheData, time.Hour)

	return results, nil
}

func (s *Storage) getCacheKeySupplier(id uuid.UUID) string {
	return fmt.Sprintf("supplier:%s", id.String())
}
