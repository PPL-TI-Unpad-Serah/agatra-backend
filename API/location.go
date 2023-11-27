package API

import (
	"agatra/model"
	service "agatra/Service" 
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
	GetLocationNearby(l *gin.Context)
	SearchLocation(l *gin.Context)
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

	var result model.LocationResponse
	result.Location = *Location 
	result.Message = "Location with ID " + strconv.Itoa(LocationID) + " Found"

	l.JSON(http.StatusOK, result)
}

func (ta *locationAPI) GetLocationList(l *gin.Context) {
	Location, err := ta.locationService.GetList()
	if err != nil {
		l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.LocationArrayResponse
	result.Locations = Location 
	result.Message = "Getting All Locations Success"

	l.JSON(http.StatusOK, result)
}

func (ta *locationAPI) GetLocationNearby(l *gin.Context) {
	Lat, err := strconv.ParseFloat(l.Query("lat"), 64)
	if err != nil {
		l.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid Lat"})
		return
	}
	Long, err := strconv.ParseFloat(l.Query("long"), 64)
	if err != nil {
		l.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid Long"})
		return
	}

	Location, err := ta.locationService.GetListNearby(Lat, Long)
	if err != nil {
		l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	var result model.LocationRangeResponse
	result.Locations = Location
	result.Message = "Location sorted by distance"
	
	l.JSON(http.StatusOK, result)
}

func (la *locationAPI) SearchLocation(l *gin.Context) {
	name := l.Query("name")
	location, err := la.locationService.SearchName(name)
	if err != nil {
		l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.LocationArrayResponse
	result.Locations = location 
	result.Message = "Getting All Privileged locations Success"

	l.JSON(http.StatusOK, result)
}