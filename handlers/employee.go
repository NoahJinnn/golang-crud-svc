package handlers

import (
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	svc  *models.EmployeeSvc
	mqtt mqtt.Client
}

func NewEmployeeHandler(svc *models.EmployeeSvc, mqtt mqtt.Client) *EmployeeHandler {
	return &EmployeeHandler{
		svc:  svc,
		mqtt: mqtt,
	}
}

// Find all employees info
// @Summary Find All Employee
// @Schemes
// @Description find all employees info
// @Produce json
// @Success 200 {array} []models.Employee
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/employees [get]
func (h *EmployeeHandler) FindAllEmployee(c *gin.Context) {
	eList, err := h.svc.FindAllEmployee(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all employees failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, eList)
}

// Find employee info by id
// @Summary Find Employee By ID
// @Schemes
// @Description find employee info by employee id
// @Produce json
// @Param        id	path	string	true	"Employee ID"
// @Success 200 {object} models.Employee
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/employee/{id} [get]
func (h *EmployeeHandler) FindEmployeeByID(c *gin.Context) {
	id := c.Param("id")

	emp, err := h.svc.FindEmployeeByID(c, id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get employee failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, emp)
}

// Create employee
// @Summary Create Employee
// @Schemes
// @Description Create employee
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagCreateEmployee	true	"Fields need to create a employee"
// @Success 200 {object} models.Employee
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/employee [post]
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	emp := &models.Employee{}
	err := c.ShouldBind(emp)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	if emp.HighestPriority {
		t := h.mqtt.Publish(mqttSvc.TOPIC_SV_HP_C, 1, false, mqttSvc.ServerUpsertHPEmployeePayload("0", emp))
		if err := mqttSvc.HandleMqttErr(&t); err != nil {
			utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        "Create HP employee mqtt failed",
				ErrorMsg:   err.Error(),
			})
			return
		}
	}

	_, err = h.svc.CreateEmployee(c.Request.Context(), emp)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create employee failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, emp)
}

// Update employee
// @Summary Update Employee By ID
// @Schemes
// @Description Update employee, must have "id" field
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagUpdateEmployee	true	"Fields need to update an employee"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/employee [patch]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	emp := &models.Employee{}
	err := c.ShouldBind(emp)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_USER_U, 1, false, mqttSvc.ServerUpdateUserPayload("0", emp.MSNV, emp.RfidPass, emp.KeypadPass, []uint{}))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update employee mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.svc.UpdateEmployee(c.Request.Context(), emp)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update employee failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

func (h *EmployeeHandler) UpdateHPEmployee(c *gin.Context) {
	emp := &models.Employee{}
	err := c.ShouldBind(emp)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_HP_U, 1, false, mqttSvc.ServerUpsertHPEmployeePayload("0", emp))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update HP employee mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.svc.UpdateHPEmployee(c.Request.Context(), emp)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update employee failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Delete employee
// @Summary Delete Employee By ID
// @Schemes
// @Description Delete employee using "id" field
// @Accept  json
// @Produce json
// @Param	data	body	object{id=int}	true	"Employee ID"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/employee [delete]
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	de := &models.DeleteEmployee{}
	err := c.ShouldBind(de)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_USER_D, 1, false, mqttSvc.ServerDeleteUserPayload("0", de.MSNV))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete employee mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.svc.DeleteEmployee(c.Request.Context(), de.MSNV)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete employee failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

func (h *EmployeeHandler) DeleteHPEmployee(c *gin.Context) {
	msnv := c.Param("msnv")

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_HP_D, 1, false, mqttSvc.ServerDeleteUserPayload("0", msnv))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete HP employee mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.svc.DeleteHPEmployee(c.Request.Context(), msnv)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete employee failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

func (h *EmployeeHandler) AppendEmployeeScheduler(c *gin.Context) {
	usu := &models.UserSchedulerUpsert{}
	empId := c.Param("id")
	err := c.ShouldBind(usu)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	emp, err := h.svc.FindEmployeeByID(c, empId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_SCHEDULER_C, 1, false, mqttSvc.ServerCreateRegisterPayload(*usu, emp.RfidPass, emp.KeypadPass, emp.MSNV))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create scheduler mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.svc.AppendEmployeeScheduler(c.Request.Context(), emp, usu.DoorlockID, &usu.Scheduler)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, true)
}

func (h *EmployeeHandler) UpdateEmployeeScheduler(c *gin.Context) {
	usu := &models.UserSchedulerUpsert{}
	empId := c.Param("id")
	err := c.ShouldBind(usu)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	emp, err := h.svc.FindEmployeeByID(c, empId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	emp, err = h.svc.UpdateEmployeeScheduler(c.Request.Context(), emp, usu.DoorlockID, &usu.Scheduler)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	scheIds := []uint{}
	for _, sche := range emp.Schedulers {
		scheIds = append(scheIds, sche.ID)
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_USER_U, 1, false,
		mqttSvc.ServerUpdateUserPayload("0", emp.MSNV, emp.RfidPass, emp.KeypadPass, scheIds))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update employee scheduler mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, true)
}

func (h *EmployeeHandler) DeleteEmployeeScheduler(c *gin.Context) {
	usu := &models.UserSchedulerUpsert{}
	empId := c.Param("id")
	err := c.ShouldBind(usu)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	emp, err := h.svc.FindEmployeeByID(c, empId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	emp, err = h.svc.DeleteEmployeeScheduler(c.Request.Context(), emp, usu.DoorlockID, &usu.Scheduler)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	scheIds := []uint{}
	for _, sche := range emp.Schedulers {
		scheIds = append(scheIds, sche.ID)
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_USER_U, 1, false,
		mqttSvc.ServerUpdateUserPayload("0", emp.MSNV, emp.RfidPass, emp.KeypadPass, scheIds))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete employee scheduler mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, true)
}
