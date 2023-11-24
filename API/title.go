package API

import (
	"agatra/model"
	service "agatra/Service" 
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TitleAPI interface {
	AddTitle(t *gin.Context)
	UpdateTitle(t *gin.Context)
	DeleteTitle(t *gin.Context)
	GetTitleByID(t *gin.Context)
	GetTitleList(t *gin.Context)
}

type titleAPI struct {
	titleService service.TitleService
}

func NewTitleAPI(titleService service.TitleService) *titleAPI{
	return &titleAPI{titleService}
}

func (ta *titleAPI) AddTitle(t *gin.Context) {
	var newTitle model.Title
	if err := t.ShouldBindJSON(&newTitle); err != nil {
		t.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	err := ta.titleService.Store(&newTitle)
	if err != nil {
		t.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	t.JSON(http.StatusOK, model.SuccessResponse{Message: "add Title success"})
}

func (ta *titleAPI) UpdateTitle(t *gin.Context) {
	var newTitle model.Title
	if err := t.ShouldBindJSON(&newTitle); err != nil {
		t.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	TitleID, err := strconv.Atoi(t.Param("id"))
	if err != nil {
		t.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Title ID"})
		return
	}
	err = ta.titleService.Update(TitleID, newTitle)
	if err != nil {
		t.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	t.JSON(http.StatusOK, model.SuccessResponse{Message: "Title update success"})
}

func (ta *titleAPI) DeleteTitle(t *gin.Context) {
	TitleID, err := strconv.Atoi(t.Param("id"))
	if err != nil {
		t.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Title ID"})
		return
	}
	err = ta.titleService.Delete(TitleID)
	if err != nil {
		t.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	t.JSON(http.StatusOK, model.SuccessResponse{Message: "Title delete success"})
}

func (ta *titleAPI) GetTitleByID(t *gin.Context) {
	titleID, err := strconv.Atoi(t.Param("id"))
	if err != nil {
		t.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid Title ID"})
		return
	}

	title, err := ta.titleService.GetByID(titleID)
	if err != nil {
		t.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.TitleResponse
	result.Title = *title 
	result.Message = "Title with ID " + strconv.Itoa(titleID) + " Found"

	t.JSON(http.StatusOK, result)
}

func (ta *titleAPI) GetTitleList(t *gin.Context) {
	Title, err := ta.titleService.GetList()
	if err != nil {
		t.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.TitleArrayResponse
	result.Titles = Title 
	result.Message = "Getting All Titles Success"

	t.JSON(http.StatusOK, result)
}