package api

import (
	"agatra/model"
	"agatra/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MachineAPI interface {
	AddMachine(m *gin.Context)
	UpdateMachine(m *gin.Context)
	DeleteMachine(m *gin.Context)
	GetMachineByID(m *gin.Context)
	GetMachineList(m *gin.Context)
}

type machineAPI struct {
	machineService service.MachineService
}

func NewMachineAPI(machineService service.MachineService) *machineAPI{
	return &machineAPI{machineService}
}

func (ma *machineAPI) AddMachine(m *gin.Context) {
	var newMachine model.Machine
	if err := m.ShouldBindJSON(&newMachine); err != nil {
		m.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	err := ma.machineService.Store(&newMachine)
	if err != nil {
		m.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	m.JSON(http.StatusOK, model.SuccessResponse{Message: "add Machine success"})
}

func (ma *machineAPI) UpdateMachine(m *gin.Context) {
	var newMachine model.Machine
	if err := m.ShouldBindJSON(&newMachine); err != nil {
		m.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}
	machineID, err := strconv.Atoi(m.Param("id"))
	if err != nil {
		m.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Machine ID"})
		return
	}
	err = ma.machineService.Update(machineID, newMachine)
	if err != nil {
		m.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	m.JSON(http.StatusOK, model.SuccessResponse{Message: "Machine update success"})
}

func (ma *machineAPI) DeleteMachine(m *gin.Context) {
	MachineID, err := strconv.Atoi(m.Param("id"))
	if err != nil {
		m.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Machine ID"})
		return
	}
	err = ma.machineService.Delete(MachineID)
	if err != nil {
		m.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}
	m.JSON(http.StatusOK, model.SuccessResponse{Message: "Machine delete success"})
}

func (ma *machineAPI) GetMachineByID(m *gin.Context) {
	machineID, err := strconv.Atoi(m.Param("id"))
	if err != nil {
		m.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid Machine ID"})
		return
	}

	machine, err := ma.machineService.GetByID(machineID)
	if err != nil {
		m.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.MachineResponse
	result.Machine = *machine 
	result.Message = "Machine with ID " + strconv.Itoa(machineID) + " Found"

	m.JSON(http.StatusOK, result)
}

func (ma *machineAPI) GetMachineList(m *gin.Context) {
	Machine, err := ma.machineService.GetList()
	if err != nil {
		m.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	var result model.MachineArrayResponse
	result.Machines = Machine 
	result.Message = "Getting All Machines Success"

	m.JSON(http.StatusOK, result)
}