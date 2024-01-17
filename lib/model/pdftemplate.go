package model

import (
	"time"
)

type PdfTemplate struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Name      string    `json:"name"`
	Elements  []Element `json:"elements" gorm:"foreignKey:PdfTemplateID"`
}

