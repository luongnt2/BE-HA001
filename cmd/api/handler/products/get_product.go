package products

import (
	"BE-HA001/pkg/httputil"
	"BE-HA001/pkg/storage"
	"errors"
	"log"
	"net/http"
)

func NewGetProductHandler(ds *storage.Storage,
	categoryStorage storage.CategoryStorage,
	supplierStorage storage.SupplierStorage) http.Handler {
	return &Product{
		ds:              ds,
		categoryStorage: categoryStorage,
		supplierStorage: supplierStorage,
	}
}

type Product struct {
	ds              *storage.Storage
	categoryStorage storage.CategoryStorage
	supplierStorage storage.SupplierStorage
}

func (h *Product) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := h.parseReq(r)

	err := h.validateReq(req)
	if err != nil {
		httputil.ResponseWrapIError(w, http.StatusBadRequest, err, 2)
		return
	}

	resp, err := getProducts(r.Context(), h.ds, h.categoryStorage, h.supplierStorage, req)
	if err != nil {
		log.Printf("Error getting products: %v", err)
		httputil.ResponseWrapIError(w, http.StatusInternalServerError, err, 1)
		return
	}

	httputil.ResponseWrapSuccessJSON(w, resp)
}

func (h *Product) parseReq(r *http.Request) storage.ProductFilter {
	return parseReqProductFilter(r)
}

func (h *Product) validateReq(req storage.ProductFilter) error {
	if req.PriceMin != nil && req.PriceMax != nil && *req.PriceMin > *req.PriceMax {
		return errors.New("invalid price")
	}
	return nil
}
