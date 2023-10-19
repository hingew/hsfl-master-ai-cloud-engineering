package pdf

import (
	"bytes"

	"github.com/jung-kurt/gofpdf"
)

type Report struct {
	pdf *gofpdf.Fpdf
}

func New() *Report {
	var r *Report

	r = &Report{}

	r.pdf = gofpdf.New("P", "mm", "A4", "")
	r.pdf.SetMargins(20, 40, 20)

	return r
}

func (r *Report) AddElement() {

}

func (r *Report) Out() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	err := r.pdf.Output(buf)

	return buf, err
}
