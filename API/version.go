package API

import (
	"agatra/model"
	service "agatra/Service" 
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VersionAPI interface {
	AddVersion(v *gin.Context)
	UpdateVersion(v *gin.Context)
	DeleteVersion(v *gin.Context)
	GetVersionByID(v *gin.Context)
	GetVersionList(v *gin.Context)
}

type versionAPI struct {
	versionService service.VersionService
}

func NewVersionAPI(versionService service.VersionService) *versionAPI{
	return &versionAPI{versionService}
}

func (va *versionAPI) AddVersion(v *gin.Context) {
	var newVersion model.Version
	if err := v.ShouldBindJSON(&newVersion); err != nil {
		v.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	err := va.versionService.Store(&newVersion)
	if err != nil {
		v.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	v.JSON(http.StatusOK, model.SuccessResponse{Message: "add Version success"})
}

func (va *versionAPI) UpdateVersion(v *gin.Context) {
	var newVersion model.Version
	if err := v.ShouldBindJSON(&newVersion); err != nil {
		v.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	VersionID, err := strconv.Atoi(v.Param("id"))
	if err != nil {
		v.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Version ID"})
		return
	}
	err = va.versionService.Update(VersionID, newVersion)
	if err != nil {
		v.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	v.JSON(http.StatusOK, model.SuccessResponse{Message: "Version update success"})
}

func (va *versionAPI) DeleteVersion(v *gin.Context) {
	VersionID, err := strconv.Atoi(v.Param("id"))
	if err != nil {
		v.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Version ID"})
		return
	}
	err = va.versionService.Delete(VersionID)
	if err != nil {
		v.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	v.JSON(http.StatusOK, model.SuccessResponse{Message: "Version delete success"})
}

func (va *versionAPI) GetVersionByID(v *gin.Context) {
	VersionID, err := strconv.Atoi(v.Param("id"))
	if err != nil {
		v.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid Version ID"})
		return
	}

	Version, err := va.versionService.GetByID(VersionID)
	if err != nil {
		v.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.VersionResponse
	result.Version = *Version 
	result.Message = "Version with ID " + strconv.Itoa(VersionID) + " Found"

	v.JSON(http.StatusOK, result)
}

func (va *versionAPI) GetVersionList(v *gin.Context) {
	titleID, err := strconv.Atoi(v.Query("game_title"))

	Version, err := va.versionService.GetList(titleID)
	if err != nil {
		v.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.VersionArrayResponse
	result.Versions = Version 
	result.Message = "Getting All Versions Success"

	v.JSON(http.StatusOK, result)
}