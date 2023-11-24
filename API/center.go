package API

import (
	"agatra/model"
	service "agatra/Service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CenterAPI interface {
	AddCenter(cr *gin.Context)
	UpdateCenter(cr *gin.Context)
	DeleteCenter(cr *gin.Context)
	GetCenterByID(cr *gin.Context)
	GetCenterList(cr *gin.Context)
}

type centerAPI struct {
	centerService service.CenterService
}

func NewCenterAPI(centerService service.CenterService) *centerAPI{
	return &centerAPI{centerService}
}

func (cra *centerAPI) AddCenter(cr *gin.Context) {
	var newCenter model.Center
	if err := cr.ShouldBindJSON(&newCenter); err != nil {
		cr.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	err := cra.centerService.Store(&newCenter)
	if err != nil {
		cr.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	cr.JSON(http.StatusOK, model.SuccessResponse{Message: "add Center success"})
}

func (cra *centerAPI) UpdateCenter(cr *gin.Context) {
	var newCenter model.Center
	if err := cr.ShouldBindJSON(&newCenter); err != nil {
		cr.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	CenterID, err := strconv.Atoi(cr.Param("id"))
	if err != nil {
		cr.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Center ID"})
		return
	}
	err = cra.centerService.Update(CenterID, newCenter)
	if err != nil {
		cr.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	cr.JSON(http.StatusOK, model.SuccessResponse{Message: "Center update success"})
}

func (cra *centerAPI) DeleteCenter(cr *gin.Context) {
	CenterID, err := strconv.Atoi(cr.Param("id"))
	if err != nil {
		cr.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Center ID"})
		return
	}
	err = cra.centerService.Delete(CenterID)
	if err != nil {
		cr.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	cr.JSON(http.StatusOK, model.SuccessResponse{Message: "Center delete success"})
}

func (cra *centerAPI) GetCenterByID(cr *gin.Context) {
	centerID, err := strconv.Atoi(cr.Param("id"))
	if err != nil {
		cr.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid Center ID"})
		return
	}

	center, err := cra.centerService.GetByID(centerID)
	if err != nil {
		cr.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.CenterResponse
	result.Center = *center 
	result.Message = "Center with ID " + strconv.Itoa(centerID) + " Found"

	cr.JSON(http.StatusOK, result)
}

func (cra *centerAPI) GetCenterList(cr *gin.Context) {
	Center, err := cra.centerService.GetList()
	if err != nil {
		cr.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.CenterArrayResponse
	result.Centers = Center 
	result.Message = "Getting All Centers Success"

	cr.JSON(http.StatusOK, result)
}