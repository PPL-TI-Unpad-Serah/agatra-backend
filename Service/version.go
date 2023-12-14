package Service

import (
	"agatra/model"
	"gorm.io/gorm"
)

type VersionService interface {
	Store(Version *model.Version) error
	Update(id int, version model.Version) error
	Delete(id int) error
	GetByID(id int) (*model.Version, error)
	GetList(titleID int) ([]model.Version, error)
}

type versionService struct {
	db *gorm.DB
}

func NewVersionService(db *gorm.DB) *versionService {
	return &versionService{db}
}

func (vs *versionService) Store(version *model.Version) error {
	return vs.db.Create(version).Error
}

func (vs *versionService) Update(id int, version model.Version) error {
	return vs.db.Where(id).Updates(version).Error
}

func (vs *versionService) Delete(id int) error {	
	return vs.db.Where(id).Delete(&model.Version{}).Error 
}

func (vs *versionService) GetByID(id int) (*model.Version, error) {
	var Version model.Version
	err := vs.db.Preload("Title").Where("id = ?", id).First(&Version).Error
	if err != nil {
		return nil, err
	}
	return &Version, nil
}

func (vs *versionService) GetList(titleID int) ([]model.Version, error) {
	var result []model.Version
	currentQuery := vs.db.Preload("Title").Order("name asc")
	if titleID != 0{
		currentQuery = currentQuery.Where("title_id = ?", titleID)
	}
	err := currentQuery.Find(&result).Error
	if err != nil{
		return []model.Version{}, err
	}
	return result, nil 
}
