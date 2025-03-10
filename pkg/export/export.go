package export

import (
	"BE-HA001/cmd/api/dto"
	"io"
)

type IExporter interface {
	Type() string
	Export(products []*dto.ListProductResponse, filename string) (IFile, string, error)
}

type IFile interface {
	Output(w io.Writer) error
}
