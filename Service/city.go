package service

import (
	"agatra/model"
	"gorm.io/gorm"
)

type CityService interface {
	Store(City *model.City) error
	Update(id int, city model.City) error
	Delete(id int) error
	GetByID(id int) (*model.City, error)
	GetList() ([]model.City, error)
}

type cityService struct {
	db *gorm.DB
}

func NewCityService(db *gorm.DB) *cityService {
	return &cityService{db}
}

func (vs *cityService) Store(city *model.City) error {
	return vs.db.Create(city).Error
}

func (vs *cityService) Update(id int, city model.City) error {
	return vs.db.Where(id).Updates(city).Error
}

func (vs *cityService) Delete(id int) error {	
	return vs.db.Where(id).Delete(&model.City{}).Error 
}

func (vs *cityService) GetByID(id int) (*model.City, error) {
	var City model.City
	err := vs.db.Where("id = ?", id).First(&City).Error
	if err != nil {
		return nil, err
	}
	return &City, nil
}

func (vs *cityService) GetList() ([]model.City, error) {
	var result []model.City
	rows, err := vs.db.Table("cities").Rows()
	if err != nil{
		return []model.City{}, err
	}
	defer rows.Close()

	for rows.Next() { 
		vs.db.ScanRows(rows, &result)
	}
	return result, nil 
}
