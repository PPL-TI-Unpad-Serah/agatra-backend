package api

import (
	"agatra/model"
	"agatra/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LocationAPI interface {
	AddLocation(l *gin.Context)
	UpdateLocation(l *gin.Context)
	DeleteLocation(l *gin.Context)
	GetLocationByID(l *gin.Context)
	GetLocationList(l *gin.Context)
}

type locationAPI struct {
	locationService service.LocationService
}

func NewLocationAPI(locationService service.LocationService) *locationAPI{
	return &locationAPI{locationService}
}

func (ta *locationAPI) AddLocation(l *gin.Context) {
	var newLocation model.Location
	if err := l.ShouldBindJSON(&newLocation); err != nil {
		l.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	err := ta.locationService.Store(&newLocation)
	if err != nil {
		l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	l.JSON(http.StatusOK, model.SuccessResponse{Message: "add Location success"})
}

func (ta *locationAPI) UpdateLocation(l *gin.Context) {
	var newLocation model.Location
	if err := l.ShouldBindJSON(&newLocation); err != nil {
		l.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	LocationID, err := strconv.Atoi(l.Param("id"))
	if err != nil {
		l.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Location ID"})
		return
	}
	err = ta.locationService.Update(LocationID, newLocation)
	if err != nil {
		l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	l.JSON(http.StatusOK, model.SuccessResponse{Message: "Location update success"})
}

func (ta *locationAPI) DeleteLocation(l *gin.Context) {
	LocationID, err := strconv.Atoi(l.Param("id"))
	if err != nil {
		l.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Location ID"})
		return
	}
	err = ta.locationService.Delete(LocationID)
	if err != nil {
		l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	l.JSON(http.StatusOK, model.SuccessResponse{Message: "Location delete success"})
}

func (ta *locationAPI) GetLocationByID(l *gin.Context) {
	LocationID, err := strconv.Atoi(l.Param("id"))
	if err != nil {
		l.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid Location ID"})
		return
	}

	Location, err := ta.locationService.GetByID(LocationID)
	if err != nil {
		l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	l.JSON(http.StatusOK, Location)
}

func (ta *locationAPI) GetLocationList(l *gin.Context) {
	Location, err := ta.locationService.GetList()
	if err != nil {
		l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	l.JSON(http.StatusOK, Location)
}