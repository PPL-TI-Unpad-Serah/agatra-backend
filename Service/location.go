package service

import (
	"agatra/model"
	// "errors"

	"gorm.io/gorm"
)

type LocationService interface {
	Store(Location *model.Location) error
	Update(id int, location model.Location) error
	Delete(id int) error
	GetByID(id int) (*model.Location, error)
	GetList() ([]model.Location, error)
	GetListNearby(lat float64, long float64) ([]model.Location_range, error)
}

type locationService struct {
	db *gorm.DB
}

func NewLocationService(db *gorm.DB) *locationService {
	return &locationService{db}
}

func (ls *locationService) Store(location *model.Location) error {
	return ls.db.Create(location).Error
}

func (ls *locationService) Update(id int, location model.Location) error {
	return ls.db.Where(id).Updates(location).Error
}

func (ls *locationService) Delete(id int) error {	
	return ls.db.Where(id).Delete(&model.Location{}).Error 
}

func (ls *locationService) GetByID(id int) (*model.Location, error) {
	var Location model.Location
	err := ls.db.Where("id = ?", id).First(&Location).Error
	if err != nil {
		return nil, err
	}
	return &Location, nil
}

func (ls *locationService) GetList() ([]model.Location, error) {
	var result []model.Location
	rows, err := ls.db.Table("locations").Rows()
	if err != nil{
		return []model.Location{}, err
	}
	defer rows.Close()

	for rows.Next() { 
		ls.db.ScanRows(rows, &result)
	}
	return result, nil 
}

func (ls *locationService) GetListNearby(lat float64, long float64) ([]model.Location_range, error) {
	var result []model.Location_range
	rows:= ls.db.Table("locations").Order("distance asc").Select("id", "name", "description", "lat", "long", gorm.Expr("(lat - ?) * (lat - ?) + (long - ?) * (long - ?) as distance", lat, lat, long, long)).Scan(&result)
	if rows.Error != nil{
		return []model.Location_range{}, rows.Error
	}
	// defer rows.Close()

	// for rows.Next() { 
	// 	ls.db.ScanRows(rows, &result)
	// }
	return result, nil 
}

