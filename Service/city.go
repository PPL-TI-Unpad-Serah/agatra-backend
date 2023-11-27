package Service

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

func (cs *cityService) Store(city *model.City) error {
	return cs.db.Create(city).Error
}

func (cs *cityService) Update(id int, city model.City) error {
	return cs.db.Where(id).Updates(city).Error
}

func (cs *cityService) Delete(id int) error {	
	return cs.db.Where(id).Delete(&model.City{}).Error 
}

func (cs *cityService) GetByID(id int) (*model.City, error) {
	var City model.City
	err := cs.db.Where("id = ?", id).First(&City).Error
	if err != nil {
		return nil, err
	}
	return &City, nil
}

func (cs *cityService) GetList() ([]model.City, error) {
	var result []model.City
	rows, err := cs.db.Table("cities").Order("name asc").Rows()
	if err != nil{
		return []model.City{}, err
	}
	defer rows.Close()

	for rows.Next() { 
		cs.db.ScanRows(rows, &result)
	}
	return result, nil 
}
