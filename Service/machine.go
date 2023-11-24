package Service

import (
	"agatra/model"
	"gorm.io/gorm"
)

type MachineService interface {
	Store(Machine *model.Machine) error
	Update(id int, machine model.Machine) error
	Delete(id int) error
	GetByID(id int) (*model.Machine, error)
	GetList() ([]model.Machine, error)
}

type machineService struct {
	db *gorm.DB
}

func NewMachineService(db *gorm.DB) *machineService {
	return &machineService{db}
}

func (ms *machineService) Store(Machine *model.Machine) error {
	return ms.db.Create(Machine).Error
}

func (ms *machineService) Update(id int, machine model.Machine) error {
	return ms.db.Where(id).Updates(machine).Error
}

func (ms *machineService) Delete(id int) error {	
	return ms.db.Where(id).Delete(&model.Machine{}).Error 
}

func (ms *machineService) GetByID(id int) (*model.Machine, error) {
	var Machine model.Machine
	err := ms.db.Where("id = ?", id).First(&Machine).Error
	if err != nil {
		return nil, err
	}
	return &Machine, nil
}

func (ms *machineService) GetList() ([]model.Machine, error) {
	var result []model.Machine
	rows, err := ms.db.Table("machines").Rows()
	if err != nil{
		return []model.Machine{}, err
	}
	defer rows.Close()

	for rows.Next() { 
		ms.db.ScanRows(rows, &result)
	}
	return result, nil 
}
