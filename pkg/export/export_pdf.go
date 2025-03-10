package export

import (
	"BE-HA001/cmd/api/dto"
	"fmt"
	"github.com/jung-kurt/gofpdf"
)

type PDF struct{}

func (p *PDF) Type() string {
	return "application/pdf"
}

func (p *PDF) Export(products []*dto.ListProductResponse, filename string) (IFile, string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()

	pdf.Cell(190, 10, "Product List")
	pdf.Ln(10)

	headers := []string{"Product ID", "Product Name", "Price", "Quantity", "Status"}
	colWidths := []float64{40, 50, 30, 20, 40}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	for _, p := range products {
		pdf.CellFormat(colWidths[0], 10, p.ProductReference, "1", 0, "", false, 0, "")
		pdf.CellFormat(colWidths[1], 10, p.ProductName, "1", 0, "", false, 0, "")
		pdf.CellFormat(colWidths[2], 10, fmt.Sprintf("%.2f", p.Price), "1", 0, "", false, 0, "")
		pdf.CellFormat(colWidths[3], 10, fmt.Sprintf("%d", p.AvailableQuantity), "1", 0, "", false, 0, "")
		pdf.CellFormat(colWidths[4], 10, p.Status, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	return pdf, fmt.Sprintf("%s.pdf", filename), nil
}
