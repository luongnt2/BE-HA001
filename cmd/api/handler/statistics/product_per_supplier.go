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

func NewProductPerSupplier(ds storage.ProductStorage,
	supplierStorage storage.SupplierStorage) *ProductPerSupplierHandler {
	return &ProductPerSupplierHandler{
		ds:              ds,
		supplierStorage: supplierStorage,
	}
}

type ProductPerSupplierHandler struct {
	ds              storage.ProductStorage
	supplierStorage storage.SupplierStorage
}

func (h *ProductPerSupplierHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	statistics, err := h.ds.CountProductBySuppliers(r.Context())
	if err != nil {
		log.Printf("failed to count product by category: %v", err)
		httputil.ResponseWrapIError(w, http.StatusInternalServerError, err, 1)
		return
	}
	if len(statistics) == 0 {
		httputil.ResponseWrapSuccessJSON(w, nil)
		return
	}

	var supplierIDs []uuid.UUID
	categoryStatisticMap := make(map[uuid.UUID]*storage.SupplierStatistic)
	for _, s := range statistics {
		supplierIDs = append(supplierIDs, s.SupplierID)
		categoryStatisticMap[s.SupplierID] = s
	}

	suppliers, err := h.supplierStorage.GetSupplierByIDs(r.Context(), supplierIDs)
	if err != nil {
		log.Printf("failed to get suppliers by ids: %v", err)
		httputil.ResponseWrapIError(w, http.StatusInternalServerError, err, 1)
		return
	}

	resp := make([]*dto.SupplierStatisticResponse, 0)
	for _, s := range suppliers {
		resp = append(resp, &dto.SupplierStatisticResponse{
			CategoryID: s.ID.String(),
			Name:       s.Name,
			Percentage: fmt.Sprintf("%.2f%%", float32(categoryStatisticMap[s.ID].ProductCount)/float32(
				total)*100),
		})
	}

	httputil.ResponseWrapSuccessJSON(w, resp)
}
