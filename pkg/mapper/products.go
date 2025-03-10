package mapper

import (
	"BE-HA001/cmd/api/dto"
	"BE-HA001/pkg/model"
	"github.com/google/uuid"
)

func ToListProductResponse(products []*model.Product, categories []*model.Category,
	suppliers []*model.Supplier) []*dto.ListProductResponse {
	var response []*dto.ListProductResponse
	categoryMap := make(map[uuid.UUID]string)
	for _, category := range categories {
		categoryMap[category.ID] = category.Name
	}
	suppliersMap := make(map[uuid.UUID]string)
	for _, supplier := range suppliers {
		suppliersMap[supplier.ID] = supplier.Name
	}

	for _, p := range products {
		response = append(response, &dto.ListProductResponse{
			ProductReference:  p.Reference,
			ProductName:       p.Name,
			DateAdded:         p.AddedDate.ToTime().Unix(),
			Status:            p.Status,
			ProductCategory:   categoryMap[p.CategoryID],
			Price:             p.Price,
			StockLocation:     p.StockCity,
			Supplier:          suppliersMap[p.SupplierID],
			AvailableQuantity: p.Quantity,
		})
	}
	return response
}
