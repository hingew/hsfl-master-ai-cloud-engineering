package model

import (
	"time"
)

type PdfTemplate struct {
	ID        int64     `json:"id"`
	PdfName   string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Elements  []Element `json:"elements"`
}

type PdfTemplateCreationRequest struct {
	PdfName  string    `json:"name"`
	Elements []Element `json:"elements"`
}

func (r *PdfTemplateCreationRequest) IsValid() bool {
	return r.PdfName != "" && r.Elements != nil
}
