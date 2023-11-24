package Service

import (
	"agatra/model"
	"gorm.io/gorm"
)

type UserService interface {
	Store(User *model.User) error
	Update(id int, user model.User) error
	Delete(id int) error
	GetByID(id int) (*model.User, error)
	GetList() ([]model.User, error)
	GetByEmail(Email string) (model.User, bool)
	GetPrivileged() ([]model.User, error)
	SearchName(name string) ([]model.User, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *userService {
	return &userService{db}
}

func (us *userService) Store(User *model.User) error {
	return us.db.Create(User).Error
}

func (us *userService) Update(id int, user model.User) error {
	return us.db.Where(id).Updates(user).Error
}

func (us *userService) Delete(id int) error {	
	return us.db.Where(id).Delete(&model.User{}).Error 
}

func (us *userService) GetByID(id int) (*model.User, error) {
	var User model.User
	err := us.db.Where("id = ?", id).First(&User).Error
	if err != nil {
		return nil, err
	}
	return &User, nil
}

func (us *userService) GetList() ([]model.User, error) {
	var result []model.User
	rows, err := us.db.Table("users").Rows()
	if err != nil{
		return []model.User{}, err
	}
	defer rows.Close()

	for rows.Next() { 
		us.db.ScanRows(rows, &result)
	}
	return result, nil 
}

func (us *userService) GetPrivileged() ([]model.User, error) {
	var result []model.User
	rows, err := us.db.Where("role = ?", "admin").Or("role = ?", "maintainer").Table("users").Rows()
	if err != nil{
		return []model.User{}, err
	}
	defer rows.Close()

	for rows.Next() { 
		us.db.ScanRows(rows, &result)
	}
	return result, nil 
}

func (us *userService) GetByEmail(email string) (model.User, bool) {
	var result model.User
	err := us.db.Where("email = ?", email).First(&result).Error
	if err != nil {
		return model.User{}, false
	}
	return result, true
}

func (us *userService) SearchName(name string) ([]model.User, error){
	var result []model.User
	rows, err := us.db.Where("name LIKE ?", "%" + name + "%").Table("users").Rows()
	if err != nil{
		return []model.User{}, err
	}
	defer rows.Close()

	for rows.Next() { 
		us.db.ScanRows(rows, &result)
	}
	return result, nil 
}
