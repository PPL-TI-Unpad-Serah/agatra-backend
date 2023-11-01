package service

import (
	"agatra/model"
	"gorm.io/gorm"
)

type CenterService interface {
	Store(Center *model.Center) error
	Update(id int, center model.Center) error
	Delete(id int) error
	GetByID(id int) (*model.Center, error)
	GetList() ([]model.Center, error)
}

type centerService struct {
	db *gorm.DB
}

func NewCenterService(db *gorm.DB) *centerService {
	return &centerService{db}
}

func (vs *centerService) Store(center *model.Center) error {
	return vs.db.Create(center).Error
}

func (vs *centerService) Update(id int, center model.Center) error {
	return vs.db.Where(id).Updates(center).Error
}

func (vs *centerService) Delete(id int) error {	
	return vs.db.Where(id).Delete(&model.Center{}).Error 
}

func (vs *centerService) GetByID(id int) (*model.Center, error) {
	var Center model.Center
	err := vs.db.Where("id = ?", id).First(&Center).Error
	if err != nil {
		return nil, err
	}
	return &Center, nil
}

func (vs *centerService) GetList() ([]model.Center, error) {
	var result []model.Center
	rows, err := vs.db.Table("centers").Rows()
	if err != nil{
		return []model.Center{}, err
	}
	defer rows.Close()

	for rows.Next() { 
		vs.db.ScanRows(rows, &result)
	}
	return result, nil 
}
