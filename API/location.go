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
}

type locationAPI struct {
	locationService service.LocationService
}

func NewLocationAPI(locationService service.LocationService) *locationAPI{
	return &locationAPI{locationService}
}

func (la *locationAPI) AddLocation(l *gin.Context) {
	var newLocation model.Location
	if err := l.ShouldBindJSON(&newLocation); err != nil {
		l.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	err := la.locationService.Store(&newLocation)
	if err != nil {
		l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	l.JSON(http.StatusOK, model.SuccessResponse{Message: "add Location success"})
}

func (la *locationAPI) UpdateLocation(l *gin.Context) {
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
	err = la.locationService.Update(LocationID, newLocation)
	if err != nil {
		l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	l.JSON(http.StatusOK, model.SuccessResponse{Message: "Location update success"})
}

func (la *locationAPI) DeleteLocation(l *gin.Context) {
	LocationID, err := strconv.Atoi(l.Param("id"))
	if err != nil {
		l.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Location ID"})
		return
	}
	err = la.locationService.Delete(LocationID)
	if err != nil {
		l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	l.JSON(http.StatusOK, model.SuccessResponse{Message: "Location delete success"})
}

func (la *locationAPI) GetLocationByID(l *gin.Context) {
	LocationID, err := strconv.Atoi(l.Param("id"))
	if err != nil {
		l.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid Location ID"})
		return
	}

	Location, err := la.locationService.GetByID(LocationID)
	if err != nil {
		l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.LocationResponse
	result.Location = *Location 
	result.Message = "Location with ID " + strconv.Itoa(LocationID) + " Found"

	l.JSON(http.StatusOK, result)
}

func (la *locationAPI) GetLocationList(l *gin.Context) {
	name := l.Query("name")
	lat, err1 := strconv.ParseFloat(l.Query("lat"), 64)
	long, err2 := strconv.ParseFloat(l.Query("long"), 64)
	pn, err := strconv.Atoi(l.Query("page"))
	if err != nil{
		pn = 1
	}
	city := l.Query("city")
	game_title_version := l.Query("game_title_version")
	game_title := l.Query("game_title")
	arcade_center := l.Query("arcade_center")

	if name != ""{
		location, err := la.locationService.SearchName(name, pn)
		if err != nil {
			l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		var result model.LocationArrayResponse
		result.Locations = location 
		result.Message = "Getting All Privileged locations Success"

		l.JSON(http.StatusOK, result)
	}else if err1 == nil && err2 == nil{
		Location, err := la.locationService.GetListNearby(lat, long, pn)
		if err != nil {
			l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}
		var result model.LocationRangeResponse
		result.Locations = Location
		result.Message = "Location sorted by distance"
		
		l.JSON(http.StatusOK, result)
	}else if city != "" || game_title_version != "" || game_title != "" || arcade_center != ""{
		Location, err := la.locationService.GetWhere(city, game_title_version, game_title, arcade_center, pn)
		if err != nil {
			l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		var result model.LocationArrayResponse
		result.Locations = Location 
		result.Message = "Getting Locations with Spec Success"

		l.JSON(http.StatusOK, result)
	}else{
		Location, err := la.locationService.GetList(pn)
		if err != nil {
			l.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		var result model.LocationArrayResponse
		result.Locations = Location 
		result.Message = "Getting All Locations Success"

		l.JSON(http.StatusOK, result)
	}
	
}

