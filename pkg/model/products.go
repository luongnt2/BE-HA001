package model

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID         uuid.UUID `gorm:"primaryKey" json:"id"`
	Reference  string    `gorm:"index" json:"reference"`
	Name       string    `json:"name"`
	AddedDate  *Date     `gorm:"type:date" json:"added_date"`
	Status     string    `json:"status"`
	CategoryID uuid.UUID `gorm:"index" json:"category_id"`
	Price      float64   `json:"price"`
	StockCity  string    `json:"stock_city"`
	SupplierID uuid.UUID `json:"supplier_id"`
	Quantity   int       `json:"available_quantity"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
