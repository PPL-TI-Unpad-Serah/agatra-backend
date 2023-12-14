package Service

import (
	"agatra/model"
	// "errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LocationService interface {
	Store(Location *model.Location) error
	Update(id int, location model.Location) error
	Delete(id int) error
	GetByID(id int) (*model.Location, error)
	GetList(page int) ([]model.Location, error)
	GetListNearby(lat float64, long float64, page int) ([]model.Location_range, error)
	SearchName(name string, page int) ([]model.Location, error)
	GetWhere(city string, version string, title string, center string, page int)([]model.Location, error)
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
	err := ls.db.Where("id = ?", id).Preload(clause.Associations).Preload("Machine.Version").Preload("Machine.Version.Title").First(&Location).Error
	if err != nil {
		return nil, err
	}
	return &Location, nil
}

func (ls *locationService) GetList(page int) ([]model.Location, error) { // TODO this should return a more compact version of Location
	var result []model.Location
	err := ls.db.Limit(10).Offset((page - 1) * 10).Preload(clause.Associations).Preload("Machine.Version").Preload("Machine.Version.Title").Find(&result).Error
	if err != nil{
		return []model.Location{}, err
	}
	return result, nil 
}

func (ls *locationService) GetListNearby(lat float64, long float64, page int) ([]model.Location_range, error) {
	var result []model.Location_range
	rows:= ls.db.Limit(10).Offset((page - 1) * 10).Table("locations").Order("distance asc").Select("id", "name", "description", "lat", "long", gorm.Expr("(lat - ?) * (lat - ?) + (long - ?) * (long - ?) as distance", lat, lat, long, long)).Scan(&result)
	if rows.Error != nil{
		return []model.Location_range{}, rows.Error
	}
	// defer rows.Close()

	// for rows.Next() { 
	// 	ls.db.ScanRows(rows, &result)
	// }
	return result, nil 
}
func (ls *locationService) SearchName(name string, page int) ([]model.Location, error){
	var result []model.Location
	rows, err := ls.db.Limit(10).Offset((page - 1) * 10).Preload(clause.Associations).Preload("Machine.Version").Preload("Machine.Version.Title").Where("name LIKE ?", "%" + name + "%").Table("locations").Rows()
	if err != nil{
		return []model.Location{}, err
	}
	defer rows.Close()

	for rows.Next() { 
		ls.db.ScanRows(rows, &result)
	}
	return result, nil 
}

func (ls *locationService) GetWhere(city string, version string, title string, center string, page int) ([]model.Location, error){
	var result []model.Location
	currentQuery := ls.db.Limit(10).Offset((page - 1) * 10)

	if city != ""{
		currentQuery = currentQuery.Where("city_id", city)
	}

	if center != ""{
		currentQuery = currentQuery.Where("center_id", center)
	}

	if version != ""{
		currentQuery = currentQuery.Where("Machine.version_id", version)
	}else if title != ""{
		currentQuery = currentQuery.Where("Machine.Version.title_id", title)
	}

	err := currentQuery.Preload(clause.Associations).Preload("Machine.Version").Preload("Machine.Version.Title").Find(&result).Error

	if err != nil{
		return []model.Location{}, err
	}
	return result, nil 
}

