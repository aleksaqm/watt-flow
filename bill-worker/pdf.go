package main

import (
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"strings"
)

func GeneratePdfFromHtml(htmlString string) ([]byte, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, fmt.Errorf("could not create PDF generator instance: %w", err)
	}

	pageReader := wkhtmltopdf.NewPageReader(strings.NewReader(htmlString))

	pdfg.AddPage(pageReader)
	err = pdfg.Create()
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF: %w", err)
	}

	pdfBytes := pdfg.Bytes()

	return pdfBytes, nil
}
