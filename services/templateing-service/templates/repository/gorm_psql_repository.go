package repository

import (
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormPsqlRepository struct {
	db *gorm.DB
}

func NewGormPsqlRepository(config database.Config) (*GormPsqlRepository, error) {
	db, err := gorm.Open(postgres.Open(config.Dsn()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &GormPsqlRepository{db}, nil
}

func (repo *GormPsqlRepository) Migrate() error {
	return repo.db.AutoMigrate(&model.PdfTemplate{}, &model.Element{} )
}

func (repo *GormPsqlRepository) CreateTemplate(data model.PdfTemplate) (*uint, error) {
	result := repo.db.Create(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	return &data.ID, nil
}

func (repo *GormPsqlRepository) GetAllTemplates() (*[]model.PdfTemplate, error) {
	var templates *[]model.PdfTemplate

	result := repo.db.Preload("Elements").Find(&templates)
	if result.Error != nil {
		return nil, result.Error
	}

	return templates, nil
}

func (repo *GormPsqlRepository) GetTemplateById(id uint) (*model.PdfTemplate, error) {
	var template model.PdfTemplate

	result := repo.db.Preload("Elements").First(&template, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &template, nil
}

func (repo *GormPsqlRepository) UpdateTemplate(id uint, updatedTemplate model.PdfTemplate) error {
	existingTemplate, err := repo.GetTemplateById(id)
	if err != nil {
		return err
	}

	if err := repo.db.Model(&existingTemplate).Association("Elements").Clear(); err != nil {
		return err
	}

	err = repo.deleteUnassignedElements()
	if err != nil {
		return err
	}

	for _, element := range updatedTemplate.Elements {
		if err := repo.db.Model(&existingTemplate).Association("Elements").Append(&element); err != nil {
			return err
		}
	}

	existingTemplate.Name = updatedTemplate.Name

	if result := repo.db.Save(&existingTemplate); result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *GormPsqlRepository) deleteUnassignedElements() error {
	var unassignedElements []model.Element
	if result := repo.db.Where("pdf_template_id IS NULL").Find(&unassignedElements); result.Error != nil {
		return result.Error
	}

	if result := repo.db.Delete(&unassignedElements); result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *GormPsqlRepository) DeleteTemplate(id uint) error {
	err := repo.db.Where("pdf_template_id = ?", id).Delete(&model.Element{}).Error
	if err != nil {
		return err
	}

	err = repo.db.Delete(&model.PdfTemplate{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
