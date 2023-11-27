package Service

import (
	"agatra/model"
	"gorm.io/gorm"
)

type TitleService interface {
	Store(Title *model.Title) error
	Update(id int, title model.Title) error
	Delete(id int) error
	GetByID(id int) (*model.Title, error)
	GetList() ([]model.Title, error)
}

type titleService struct {
	db *gorm.DB
}

func NewTitleService(db *gorm.DB) *titleService {
	return &titleService{db}
}

func (ts *titleService) Store(title *model.Title) error {
	return ts.db.Create(title).Error
}

func (ts *titleService) Update(id int, title model.Title) error {
	return ts.db.Where(id).Updates(title).Error
}

func (ts *titleService) Delete(id int) error {	
	return ts.db.Where(id).Delete(&model.Title{}).Error 
}

func (ts *titleService) GetByID(id int) (*model.Title, error) {
	var Title model.Title
	err := ts.db.Where("id = ?", id).First(&Title).Error
	if err != nil {
		return nil, err
	}
	return &Title, nil
}

func (ts *titleService) GetList() ([]model.Title, error) {
	var result []model.Title
	rows, err := ts.db.Table("titles").Order("name asc").Rows()
	if err != nil{
		return []model.Title{}, err
	}
	defer rows.Close()

	for rows.Next() { 
		ts.db.ScanRows(rows, &result)
	}
	return result, nil 
}
