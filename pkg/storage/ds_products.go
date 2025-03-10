package storage

import (
	"BE-HA001/pkg/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"strings"
	"time"
)

type ProductStorage interface {
	GetProducts(ctx context.Context, filter ProductFilter) ([]*model.Product, error)
	CountProductByCategory(ctx context.Context) ([]*CategoryStatistic, error)
	CountProductBySuppliers(ctx context.Context) ([]*SupplierStatistic, error)
	CountTotalProduct(ctx context.Context) (int64, error)
}

type (
	ProductFilter struct {
		Name            *string    `json:"name" query:"name ILIKE ?"`
		Reference       *string    `json:"reference" query:"reference = ?"`
		Status          []string   `json:"status" query:"status in ?"`
		Category        []string   `json:"category" query:"category_id in ?"`
		StockCity       []string   `json:"stock_city" query:"stock_city in ?"`
		Supplier        *string    `json:"supplier" query:"supplier_id = ?"`
		PriceMin        *float64   `json:"price_min" query:"price >= ?"`
		PriceMax        *float64   `json:"price_max" query:"price <= ?"`
		AvailableMin    *int       `json:"available_min" query:"quantity >= ?"`
		AvailableMax    *int       `json:"available_max" query:"quantity <= ?"`
		DateFrom        *time.Time `json:"date_from" query:"added_date >= ?"`
		DateTo          *time.Time `json:"date_to" query:"added_date <= ?"`
		BeforeCreatedAt *time.Time `json:"before_created_at" query:"created_at > ?"`
		Limit           int        `json:"limit" query:"-"`
	}

	CategoryStatistic struct {
		CategoryID   uuid.UUID `json:"category_id"`
		ProductCount int64     `json:"product_count"`
	}

	SupplierStatistic struct {
		SupplierID   uuid.UUID `json:"supplier_id"`
		ProductCount int64     `json:"product_count"`
	}
)

func (s *Storage) GetProducts(ctx context.Context, filter ProductFilter) ([]*model.Product, error) {
	query := s.DB.Model(&model.Product{})

	v := reflect.ValueOf(filter)
	t := reflect.TypeOf(filter)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		sqlCond := field.Tag.Get("query")
		if sqlCond == "" || sqlCond == "-" {
			continue
		}

		if value.Kind() == reflect.Ptr && !value.IsNil() {
			if strings.Contains(sqlCond, "ILIKE") {
				query = query.Where(sqlCond, fmt.Sprintf("%s%%", value.Elem().Interface()))
			} else {
				query = query.Where(sqlCond, value.Elem().Interface())
			}
		}

		if value.Kind() == reflect.Slice && value.Len() > 0 {
			query = query.Where(sqlCond, value.Interface())
		}
	}

	var products []*model.Product
	if err := query.WithContext(ctx).
		Order("created_at DESC").
		Limit(filter.Limit).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Storage) CountProductByCategory(ctx context.Context) ([]*CategoryStatistic, error) {
	categories := make([]*CategoryStatistic, 0)
	if err := s.WithContext(ctx).Model(&model.Product{}).
		Select("category_id, count(*) as product_count").
		Group("category_id").
		Scan(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *Storage) CountProductBySuppliers(ctx context.Context) ([]*SupplierStatistic, error) {
	resp := make([]*SupplierStatistic, 0)
	if err := s.WithContext(ctx).Model(&model.Product{}).
		Select("supplier_id, count(*) as product_count").
		Group("supplier_id").
		Scan(&resp).Error; err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Storage) CountTotalProduct(ctx context.Context) (int64, error) {
	var totalCount int64
	if err := s.WithContext(ctx).Model(&model.Product{}).Count(&totalCount).Error; err != nil {
		return 0, err
	}

	return totalCount, nil
}
