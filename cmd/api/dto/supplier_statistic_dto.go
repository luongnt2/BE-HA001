package dto

type SupplierStatisticResponse struct {
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
	Percentage string `json:"percentage"`
}
