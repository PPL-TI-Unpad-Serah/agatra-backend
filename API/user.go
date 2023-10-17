package api

import (
	"agatra/model"
	"agatra/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserAPI interface {
	AddUser(u *gin.Context)
	UpdateUser(u *gin.Context)
	DeleteUser(u *gin.Context)
	GetUserByID(u *gin.Context)
	GetUserList(u *gin.Context)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI{
	return &userAPI{userService}
}


func (ua *userAPI) AddUser(m *gin.Context) {
	var newUser model.User
	if err := m.ShouldBindJSON(&newUser); err != nil {
		m.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	err := ua.userService.Store(&newUser)
	if err != nil {
		m.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	m.JSON(http.StatusOK, model.SuccessResponse{Message: "add User success"})
}

func (ua *userAPI) UpdateUser(m *gin.Context) {
	var newUser model.User
	if err := m.ShouldBindJSON(&newUser); err != nil {
		m.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	UserID, err := strconv.Atoi(m.Param("id"))
	if err != nil {
		m.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid User ID"})
		return
	}
	err = ua.userService.Update(UserID, newUser)
	if err != nil {
		m.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	m.JSON(http.StatusOK, model.SuccessResponse{Message: "User update success"})
}

func (ua *userAPI) DeleteUser(m *gin.Context) {
	UserID, err := strconv.Atoi(m.Param("id"))
	if err != nil {
		m.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid User ID"})
		return
	}
	err = ua.userService.Delete(UserID)
	if err != nil {
		m.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	m.JSON(http.StatusOK, model.SuccessResponse{Message: "User delete success"})
}

func (ua *userAPI) GetUserByID(m *gin.Context) {
	UserID, err := strconv.Atoi(m.Param("id"))
	if err != nil {
		m.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid User ID"})
		return
	}

	User, err := ua.userService.GetByID(UserID)
	if err != nil {
		m.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	m.JSON(http.StatusOK, User)
}

func (ua *userAPI) GetUserList(m *gin.Context) {// TODO: answer here
	User, err := ua.userService.GetList()
	if err != nil {
		m.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	m.JSON(http.StatusOK, User)
}