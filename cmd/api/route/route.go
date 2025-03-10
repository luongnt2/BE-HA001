package router

import (
	"BE-HA001/cmd/api/handler/locations"
	"BE-HA001/cmd/api/handler/products"
	"BE-HA001/cmd/api/handler/statistics"
	"BE-HA001/cmd/api/middleware"
	"BE-HA001/pkg/export"
	"BE-HA001/pkg/storage"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(ds *storage.Storage, exportHelper export.IExporter) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/products", products.NewGetProductHandler(ds, ds, ds).ServeHTTP).Methods(http.MethodGet)
	r.HandleFunc("/products/export", products.NewExportProductsHandler(ds, ds, ds,
		exportHelper).ServeHTTP).Methods(http.MethodPost)
	r.HandleFunc("/distance", middleware.IPMiddleware(
		locations.NewGetDistanceHandler()).ServeHTTP).Methods(http.MethodGet)
	r.HandleFunc("/api/statistics/products-per-category",
		statistics.NewProductPerCategory(ds, ds).ServeHTTP).Methods(http.MethodGet)
	r.HandleFunc("/api/statistics/products-per-supplier",
		statistics.NewProductPerSupplier(ds, ds).ServeHTTP).Methods(http.MethodGet)

	return r
}
