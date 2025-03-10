package statistics

import (
	"BE-HA001/cmd/api/dto"
	"BE-HA001/pkg/httputil"
	"BE-HA001/pkg/storage"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func NewProductPerCategory(ds storage.ProductStorage, categoryStorage storage.CategoryStorage) *ProductPerCategoryHandler {
	return &ProductPerCategoryHandler{
		ds:              ds,
		categoryStorage: categoryStorage,
	}
}

type ProductPerCategoryHandler struct {
	ds              storage.ProductStorage
	categoryStorage storage.CategoryStorage
}

func (h *ProductPerCategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	total, err := h.ds.CountTotalProduct(r.Context())
	if err != nil {
		log.Printf("failed to count total product: %v", err)
		httputil.ResponseWrapIError(w, http.StatusInternalServerError, err, 1)
		return
	}
	if total == 0 {
		httputil.ResponseWrapSuccessJSON(w, nil)
		return
	}

	categoryStatistics, err := h.ds.CountProductByCategory(r.Context())
	if err != nil {
		log.Printf("failed to count product by category: %v", err)
		httputil.ResponseWrapIError(w, http.StatusInternalServerError, err, 1)
		return
	}
	if len(categoryStatistics) == 0 {
		httputil.ResponseWrapSuccessJSON(w, nil)
		return
	}

	var categoryIDs []uuid.UUID
	categoryStatisticMap := make(map[uuid.UUID]*storage.CategoryStatistic)
	for _, category := range categoryStatistics {
		categoryIDs = append(categoryIDs, category.CategoryID)
		categoryStatisticMap[category.CategoryID] = category
	}

	categories, err := h.categoryStorage.GetCategoriesByIDs(r.Context(), categoryIDs)
	if err != nil {
		log.Printf("failed to get categories by ids: %v", err)
		httputil.ResponseWrapIError(w, http.StatusInternalServerError, err, 1)
		return
	}

	resp := make([]*dto.CategoryStatisticResponse, 0)
	for _, category := range categories {
		resp = append(resp, &dto.CategoryStatisticResponse{
			CategoryID: category.ID.String(),
			Name:       category.Name,
			Percentage: fmt.Sprintf("%.2f%%", float32(categoryStatisticMap[category.ID].ProductCount)/float32(
				total)*100),
		})
	}

	httputil.ResponseWrapSuccessJSON(w, resp)
}
