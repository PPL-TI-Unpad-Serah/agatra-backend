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
	GetList() ([]model.Version, error)
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
	err := vs.db.Preload("titles").Where("id = ?", id).First(&Version).Error
	if err != nil {
		return nil, err
	}
	return &Version, nil
}

func (vs *versionService) GetList() ([]model.Version, error) {
	var result []model.Version
	rows, err := vs.db.Preload("titles").Table("versions").Rows()
	if err != nil{
		return []model.Version{}, err
	}
	defer rows.Close()

	for rows.Next() { 
		vs.db.ScanRows(rows, &result)
	}
	return result, nil 
}
