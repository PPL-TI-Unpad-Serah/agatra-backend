package api

import (
	"agatra/model"
	"agatra/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CentreAPI interface {
	AddCentre(cr *gin.Context)
	UpdateCentre(cr *gin.Context)
	DeleteCentre(cr *gin.Context)
	GetCentreByID(cr *gin.Context)
	GetCentreList(cr *gin.Context)
}

type centreAPI struct {
	centreService service.CentreService
}

func NewCentreAPI(centreService service.CentreService) *centreAPI{
	return &centreAPI{centreService}
}

func (cra *centreAPI) AddCentre(cr *gin.Context) {
	var newCentre model.Centre
	if err := cr.ShouldBindJSON(&newCentre); err != nil {
		cr.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	err := cra.centreService.Store(&newCentre)
	if err != nil {
		cr.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	cr.JSON(http.StatusOK, model.SuccessResponse{Message: "add Centre success"})
}

func (cra *centreAPI) UpdateCentre(cr *gin.Context) {
	var newCentre model.Centre
	if err := cr.ShouldBindJSON(&newCentre); err != nil {
		cr.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	CentreID, err := strconv.Atoi(cr.Param("id"))
	if err != nil {
		cr.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Centre ID"})
		return
	}
	err = cra.centreService.Update(CentreID, newCentre)
	if err != nil {
		cr.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	cr.JSON(http.StatusOK, model.SuccessResponse{Message: "Centre update success"})
}

func (cra *centreAPI) DeleteCentre(cr *gin.Context) {
	CentreID, err := strconv.Atoi(cr.Param("id"))
	if err != nil {
		cr.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Centre ID"})
		return
	}
	err = cra.centreService.Delete(CentreID)
	if err != nil {
		cr.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	cr.JSON(http.StatusOK, model.SuccessResponse{Message: "Centre delete success"})
}

func (cra *centreAPI) GetCentreByID(cr *gin.Context) {
	CentreID, err := strconv.Atoi(cr.Param("id"))
	if err != nil {
		cr.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid Centre ID"})
		return
	}

	Centre, err := cra.centreService.GetByID(CentreID)
	if err != nil {
		cr.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	cr.JSON(http.StatusOK, Centre)
}

func (cra *centreAPI) GetCentreList(cr *gin.Context) {
	Centre, err := cra.centreService.GetList()
	if err != nil {
		cr.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	cr.JSON(http.StatusOK, Centre)
}