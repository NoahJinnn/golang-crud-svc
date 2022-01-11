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
	svc     *models.EmployeeSvc
	scheSvc *models.SchedulerSvc
	mqtt    mqtt.Client
}

func NewEmployeeHandler(svc *models.EmployeeSvc, scheSvc *models.SchedulerSvc, mqtt mqtt.Client) *EmployeeHandler {
	return &EmployeeHandler{
		svc:     svc,
		scheSvc: scheSvc,
		mqtt:    mqtt,
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

// Find employee info by msnv
// @Summary Find Employee By MSNV
// @Schemes
// @Description find employee info by employee msnv
// @Produce json
// @Param        msnv	path	string	true	"Employee MSNV"
// @Success 200 {object} models.Employee
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/employee/{msnv} [get]
func (h *EmployeeHandler) FindEmployeeByMSNV(c *gin.Context) {
	msnv := c.Param("msnv")

	emp, err := h.svc.FindEmployeeByMSNV(c, msnv)
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
		t := h.mqtt.Publish(mqttSvc.TOPIC_SV_HP_C, 1, false,
			mqttSvc.ServerUpdateUserPayload("0", emp.MSNV, emp.RfidPass, emp.KeypadPass))
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
// @Summary Update Employee By ID and MSNV
// @Schemes
// @Description Update employee, must have "id" and "msnv" field
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

	isSuccess, err := h.svc.UpdateEmployee(c.Request.Context(), emp)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update employee failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_USER_U, 1, false,
		mqttSvc.ServerUpdateUserPayload("0", emp.MSNV, emp.RfidPass, emp.KeypadPass))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update employee mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	if emp.HighestPriority {
		t = h.mqtt.Publish(mqttSvc.TOPIC_SV_USER_U, 1, false,
			mqttSvc.ServerUpdateUserPayload("0", emp.MSNV, emp.RfidPass, emp.KeypadPass))
		if err := mqttSvc.HandleMqttErr(&t); err != nil {
			utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        "Update employee mqtt failed",
				ErrorMsg:   err.Error(),
			})
			return
		}
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

	isSuccess, err := h.svc.UpdateEmployee(c.Request.Context(), emp)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update employee failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	if emp.HighestPriority {
		t := h.mqtt.Publish(mqttSvc.TOPIC_SV_HP_U, 1, false,
			mqttSvc.ServerUpdateUserPayload("0", emp.MSNV, emp.RfidPass, emp.KeypadPass))
		if err := mqttSvc.HandleMqttErr(&t); err != nil {
			utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        "Update HP employee mqtt failed",
				ErrorMsg:   err.Error(),
			})
			return
		}
	} else {
		t := h.mqtt.Publish(mqttSvc.TOPIC_SV_HP_D, 1, false,
			mqttSvc.ServerDeleteUserPayload("0", emp.MSNV))
		if err := mqttSvc.HandleMqttErr(&t); err != nil {
			utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Msg:        "Update HP employee mqtt failed",
				ErrorMsg:   err.Error(),
			})
			return
		}
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Delete employee
// @Summary Delete Employee By ID
// @Schemes
// @Description Delete employee using "msnv" field
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

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_HP_D, 1, false,
		mqttSvc.ServerDeleteUserPayload("0", msnv))
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
	msnv := c.Param("msnv")
	err := c.ShouldBind(usu)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	emp, err := h.svc.FindEmployeeByMSNV(c, msnv)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get employee failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	sche := &usu.Scheduler
	sche.EmployeeID = &emp.MSNV
	sche.DoorSerialID = &usu.DoorlockID
	_, err = h.scheSvc.CreateScheduler(c.Request.Context(), sche)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create scheduler failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_SCHEDULER_C, 1, false, mqttSvc.ServerCreateRegisterPayload(
		usu.GatewayID,
		usu.DoorlockID,
		sche,
		&mqttSvc.UserIDPassword{
			UserId:     emp.MSNV,
			RfidPass:   emp.RfidPass,
			KeypadPass: emp.KeypadPass,
		}))

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
			Msg:        "Update employee failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, true)
}

// func (h *EmployeeHandler) UpdateEmployeeScheduler(c *gin.Context) {
// 	usu := &models.UserSchedulerUpsert{}
// 	msnv := c.Param("msnv")
// 	err := c.ShouldBind(usu)
// 	if err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Invalid req body",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	emp, err := h.svc.FindEmployeeByMSNV(c, msnv)
// 	if err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Get employee failed",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	emp, err = h.svc.UpdateEmployeeScheduler(c.Request.Context(), emp, usu.DoorlockID, &usu.Scheduler)
// 	if err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Update employee failed",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	scheIds := []uint{}
// 	for _, sche := range emp.Schedulers {
// 		scheIds = append(scheIds, sche.ID)
// 	}

// 	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_USER_U, 1, false,
// 		mqttSvc.ServerUpdateUserPayload("0", emp.MSNV, emp.RfidPass, emp.KeypadPass, scheIds))
// 	if err := mqttSvc.HandleMqttErr(&t); err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Update employee scheduler mqtt failed",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	utils.ResponseJson(c, http.StatusOK, true)
// }

// func (h *EmployeeHandler) DeleteEmployeeScheduler(c *gin.Context) {
// 	usu := &models.UserSchedulerUpsert{}
// 	msnv := c.Param("msnv")
// 	err := c.ShouldBind(usu)
// 	if err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Invalid req body",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	emp, err := h.svc.FindEmployeeByMSNV(c, msnv)
// 	if err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Get employee failed",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	emp, err = h.svc.DeleteEmployeeScheduler(c.Request.Context(), emp, usu.DoorlockID, &usu.Scheduler)
// 	if err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Update employee failed",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	scheIds := []uint{}
// 	for _, sche := range emp.Schedulers {
// 		scheIds = append(scheIds, sche.ID)
// 	}

// 	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_USER_U, 1, false,
// 		mqttSvc.ServerUpdateUserPayload("0", emp.MSNV, emp.RfidPass, emp.KeypadPass, scheIds))
// 	if err := mqttSvc.HandleMqttErr(&t); err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Delete employee scheduler mqtt failed",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	utils.ResponseJson(c, http.StatusOK, true)
// }
