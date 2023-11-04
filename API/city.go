package api

import (
	"agatra/model"
	"agatra/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CityAPI interface {
	AddCity(cy *gin.Context)
	UpdateCity(cy *gin.Context)
	DeleteCity(cy *gin.Context)
	GetCityByID(cy *gin.Context)
	GetCityList(cy *gin.Context)
}

type cityAPI struct {
	cityService service.CityService
}

func NewCityAPI(cityService service.CityService) *cityAPI{
	return &cityAPI{cityService}
}

func (cya *cityAPI) AddCity(cy *gin.Context) {
	var newCity model.City
	if err := cy.ShouldBindJSON(&newCity); err != nil {
		cy.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	err := cya.cityService.Store(&newCity)
	if err != nil {
		cy.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	cy.JSON(http.StatusOK, model.SuccessResponse{Message: "add City success"})
}

func (cya *cityAPI) UpdateCity(cy *gin.Context) {
	var newCity model.City
	if err := cy.ShouldBindJSON(&newCity); err != nil {
		cy.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	CityID, err := strconv.Atoi(cy.Param("id"))
	if err != nil {
		cy.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid City ID"})
		return
	}
	err = cya.cityService.Update(CityID, newCity)
	if err != nil {
		cy.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	cy.JSON(http.StatusOK, model.SuccessResponse{Message: "City update success"})
}

func (cya *cityAPI) DeleteCity(cy *gin.Context) {
	CityID, err := strconv.Atoi(cy.Param("id"))
	if err != nil {
		cy.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid City ID"})
		return
	}
	err = cya.cityService.Delete(CityID)
	if err != nil {
		cy.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	cy.JSON(http.StatusOK, model.SuccessResponse{Message: "City delete success"})
}

func (cya *cityAPI) GetCityByID(cy *gin.Context) {
	cityID, err := strconv.Atoi(cy.Param("id"))
	if err != nil {
		cy.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid City ID"})
		return
	}

	city, err := cya.cityService.GetByID(cityID)
	if err != nil {
		cy.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.CityResponse
	result.City = *city 
	result.Message = "City with ID " + strconv.Itoa(cityID) + " Found"

	cy.JSON(http.StatusOK, result)
}

func (cya *cityAPI) GetCityList(cy *gin.Context) {
	City, err := cya.cityService.GetList()
	if err != nil {
		cy.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.CityArrayResponse
	result.Cities = City 
	result.Message = "Getting All Cities Success"

	cy.JSON(http.StatusOK, result)
}