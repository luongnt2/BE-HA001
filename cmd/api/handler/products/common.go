package products

import (
	"BE-HA001/cmd/api/dto"
	"BE-HA001/pkg/mapper"
	"BE-HA001/pkg/storage"
	"context"
	"github.com/google/uuid"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func parseReqProductFilter(r *http.Request) storage.ProductFilter {
	filter := storage.ProductFilter{}

	v := reflect.ValueOf(&filter).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		queryKey := field.Tag.Get("json")

		if queryKey == "" {
			continue
		}

		param := r.URL.Query().Get(queryKey)
		if param == "" {
			continue
		}
		if v.Field(i).Type() == reflect.TypeOf(&time.Time{}) {
			if parsed, err := strconv.ParseInt(param, 10, 64); err == nil {
				timeUnix := time.Unix(parsed, 0)
				v.Field(i).Set(reflect.ValueOf(&timeUnix))
			}
			continue
		}

		switch v.Field(i).Kind() {
		case reflect.String:
			v.Field(i).SetString(param)
		case reflect.Slice:
			v.Field(i).Set(reflect.ValueOf(strings.Split(param, ",")))
		case reflect.Ptr:
			switch field.Type.Elem().Kind() {
			case reflect.String:
				v.Field(i).Set(reflect.ValueOf(&param))
			case reflect.Int, reflect.Int64:
				if parsed, err := strconv.ParseInt(param, 10, 64); err == nil {
					v.Field(i).Set(reflect.ValueOf(&parsed))
				}
			case reflect.Float64:
				if parsed, err := strconv.ParseFloat(param, 64); err == nil {
					v.Field(i).Set(reflect.ValueOf(&parsed))
				}
			default:
				panic("unhandled default case")
			}
		case reflect.Int:
			if parsed, err := strconv.Atoi(param); err == nil {
				v.Field(i).SetInt(int64(parsed))
			}
		case reflect.Int64:
			if parsed, err := strconv.ParseInt(param, 10, 64); err == nil {
				v.Field(i).SetInt(parsed)
			}
		case reflect.Float64:
			if parsed, err := strconv.ParseFloat(param, 64); err == nil {
				v.Field(i).SetFloat(parsed)
			}
		default:
			panic("unhandled default case")
		}
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	return filter
}

func getProducts(ctx context.Context,
	prodStorage storage.ProductStorage,
	categoryStorage storage.CategoryStorage,
	supplierStorage storage.SupplierStorage,
	req storage.ProductFilter) (
	[]*dto.ListProductResponse, error) {
	products, err := prodStorage.GetProducts(ctx, req)
	if err != nil {
		log.Printf("Error getting products: %v", err)
		return nil, err
	}

	categoryIds := make([]uuid.UUID, 0)
	supplierIds := make([]uuid.UUID, 0)
	for _, product := range products {
		categoryIds = append(categoryIds, product.CategoryID)
		supplierIds = append(supplierIds, product.SupplierID)
	}

	categories, err := categoryStorage.GetCategoriesByIDs(ctx, categoryIds)
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		return nil, err
	}

	suppliers, err := supplierStorage.GetSupplierByIDs(ctx, supplierIds)
	if err != nil {
		log.Printf("Error getting suppliers: %v", err)
		return nil, err
	}

	return mapper.ToListProductResponse(products, categories, suppliers), nil
}
