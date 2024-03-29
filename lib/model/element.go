package model

type Element struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Type          string `json:"type"`
	X             int    `json:"x"`
	Y             int    `json:"y"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	ValueFrom     string `json:"value_from"`
	Font          string `json:"font"`
	FontSize      int    `json:"font_size"`
	PdfTemplateID uint   `json:"pdf_template_id"`
}
