package products

import (
	"BE-HA001/pkg/export"
	"BE-HA001/pkg/httputil"
	"BE-HA001/pkg/storage"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ExportProductsHandler struct {
	ds              storage.ProductStorage
	categoryStorage storage.CategoryStorage
	supplierStorage storage.SupplierStorage
	exportHelper    export.IExporter
}

func NewExportProductsHandler(
	ds storage.ProductStorage,
	categoryStorage storage.CategoryStorage,
	supplierStorage storage.SupplierStorage,
	exportHelper export.IExporter) *ExportProductsHandler {
	return &ExportProductsHandler{
		ds:              ds,
		categoryStorage: categoryStorage,
		supplierStorage: supplierStorage,
		exportHelper:    exportHelper,
	}
}

func (h *ExportProductsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := parseReqProductFilter(r)

	// Set max export, Avoid leaking memory when the number of files is too large
	req.Limit = 1000
	req.BeforeCreatedAt = nil

	products, err := getProducts(r.Context(), h.ds, h.categoryStorage, h.supplierStorage, req)
	if err != nil {
		log.Printf("Error getting products: %v", err)
		httputil.ResponseWrapIError(w, http.StatusInternalServerError, err, 1)
		return
	}

	file, fileName, err := h.exportHelper.Export(products,
		fmt.Sprintf("product_%s", time.Now().Format("20060201")))
	if err != nil {
		log.Printf("Error exporting products: %v", err)
		httputil.ResponseWrapIError(w, http.StatusInternalServerError, err, 1)
		return
	}

	w.Header().Set("Content-Type", h.exportHelper.Type())
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	err = file.Output(w)
	if err != nil {
		httputil.ResponseWrapIError(w, http.StatusInternalServerError, err, 1)
		return
	}

	httputil.ResponseWrapSuccessJSON(w, nil)
}
