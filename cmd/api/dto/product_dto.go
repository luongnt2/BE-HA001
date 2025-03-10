package dto

type ListProductResponse struct {
	ProductReference  string  `json:"product_reference"`
	ProductName       string  `json:"product_name"`
	DateAdded         int64   `json:"date_added"`
	Status            string  `json:"status"`
	ProductCategory   string  `json:"product_category"`
	Price             float64 `json:"price"`
	StockLocation     string  `json:"stock_location"`
	Supplier          string  `json:"supplier"`
	AvailableQuantity int     `json:"available_quantity"`
}
