package pdf

import (
	"bytes"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
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

func (r *Report) Render(template *model.PdfTemplate, params map[string]interface{}) {
	for _, el := range template.Elements {
		r.renderElement(el, params)
	}

	r.pdf.SetTitle(template.Name, true)
}

func (r *Report) renderElement(element model.Element, params map[string]interface{}) {
	switch element.Type {
	case "rect":
		r.pdf.Rect(float64(element.X), float64(element.Y), float64(element.Width), float64(element.Height), "")

	case "text":
		r.pdf.SetFont(element.Font, "", float64(element.FontSize))
		r.pdf.Text(float64(element.X), float64(element.Y), params[element.ValueFrom].(string))
	}

}

func (r *Report) Out() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	err := r.pdf.Output(buf)

	return buf, err
}
