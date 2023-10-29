package service

import (
	"agatra/model"
	"gorm.io/gorm"
)

type CentreService interface {
	Store(Centre *model.Centre) error
	Update(id int, centre model.Centre) error
	Delete(id int) error
	GetByID(id int) (*model.Centre, error)
	GetList() ([]model.Centre, error)
}

type centreService struct {
	db *gorm.DB
}

func NewCentreService(db *gorm.DB) *centreService {
	return &centreService{db}
}

func (vs *centreService) Store(centre *model.Centre) error {
	return vs.db.Create(centre).Error
}

func (vs *centreService) Update(id int, centre model.Centre) error {
	return vs.db.Where(id).Updates(centre).Error
}

func (vs *centreService) Delete(id int) error {	
	return vs.db.Where(id).Delete(&model.Centre{}).Error 
}

func (vs *centreService) GetByID(id int) (*model.Centre, error) {
	var Centre model.Centre
	err := vs.db.Where("id = ?", id).First(&Centre).Error
	if err != nil {
		return nil, err
	}
	return &Centre, nil
}

func (vs *centreService) GetList() ([]model.Centre, error) {
	var result []model.Centre
	rows, err := vs.db.Table("centres").Rows()
	if err != nil{
		return []model.Centre{}, err
	}
	defer rows.Close()

	for rows.Next() { 
		vs.db.ScanRows(rows, &result)
	}
	return result, nil 
}
