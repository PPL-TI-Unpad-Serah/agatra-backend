package Service

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

func (cs *centerService) Store(center *model.Center) error {
	return cs.db.Create(center).Error
}

func (cs *centerService) Update(id int, center model.Center) error {
	return cs.db.Where(id).Updates(center).Error
}

func (cs *centerService) Delete(id int) error {	
	return cs.db.Where(id).Delete(&model.Center{}).Error 
}

func (cs *centerService) GetByID(id int) (*model.Center, error) {
	var Center model.Center
	err := cs.db.Where("id = ?", id).First(&Center).Error
	if err != nil {
		return nil, err
	}
	return &Center, nil
}

func (cs *centerService) GetList() ([]model.Center, error) {
	var result []model.Center
	rows, err := cs.db.Table("centers").Rows()
	if err != nil{
		return []model.Center{}, err
	}
	defer rows.Close()

	for rows.Next() { 
		cs.db.ScanRows(rows, &result)
	}
	return result, nil 
}
