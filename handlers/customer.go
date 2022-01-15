package handlers

import (
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	svc     *models.CustomerSvc
	scheSvc *models.SchedulerSvc
	mqtt    mqtt.Client
}

func NewCustomerHandler(svc *models.CustomerSvc, scheSvc *models.SchedulerSvc, mqtt mqtt.Client) *CustomerHandler {
	return &CustomerHandler{
		svc:     svc,
		scheSvc: scheSvc,
		mqtt:    mqtt,
	}
}

// Find all customers info
// @Summary Find All Customer
// @Schemes
// @Description find all customers info
// @Produce json
// @Success 200 {array} []models.Customer
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/customers [get]
func (h *CustomerHandler) FindAllCustomer(c *gin.Context) {
	sList, err := h.svc.FindAllCustomer(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all customers failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, sList)
}

// Find customer info by cccd
// @Summary Find Customer By CCCD
// @Schemes
// @Description find customer info by customer cccd
// @Produce json
// @Param        cccd	path	string	true	"Customer CCCD"
// @Success 200 {object} models.Customer
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/customer/{cccd} [get]
func (h *CustomerHandler) FindCustomerByCCCD(c *gin.Context) {
	cccd := c.Param("cccd")

	cus, err := h.svc.FindCustomerByCCCD(c, cccd)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get customer failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, cus)
}

// Create customer
// @Summary Create Customer
// @Schemes
// @Description Create customer
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagCreateCustomer	true	"Fields need to create a customer"
// @Success 200 {object} models.Customer
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/customer [post]
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	cus := &models.Customer{}
	err := c.ShouldBind(cus)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.svc.CreateCustomer(c.Request.Context(), cus)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create customer failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, cus)
}

// Update customer
// @Summary Update Customer By ID
// @Schemes
// @Description Update customer, must have correct "id" and "cccd" field
// @Accept  json
// @Produce json
// @Param	data	body	models.Customer	true	"Fields need to update a customer"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/customer [patch]
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	cus := &models.Customer{}
	err := c.ShouldBind(cus)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.svc.UpdateCustomer(c.Request.Context(), cus)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update customer failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_USER_U, 1, false,
		mqttSvc.ServerUpdateUserPayload("0", cus.CCCD, cus.RfidPass, cus.KeypadPass))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update customer mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Delete customer
// @Summary Delete Customer By CCCD
// @Schemes
// @Description Delete customer using "cccd" field
// @Accept  json
// @Produce json
// @Param	data	body	models.DeleteCustomer	true	"Customer CCCD"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/customer [delete]
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	dcus := &models.DeleteCustomer{}
	err := c.ShouldBind(dcus)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.svc.DeleteCustomer(c.Request.Context(), dcus.CCCD)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete customer failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_USER_D, 1, false,
		mqttSvc.ServerDeleteUserPayload("0", dcus.CCCD))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete customer mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Add customer scheduler
// @Summary Add Door Open Scheduler For Customer
// @Schemes
// @Description Add scheduler that allows customer open specific door
// @Accept  json
// @Produce json
// @Param	data	body	models.UserSchedulerUpsert	true	"Request with Scheduler, GatewayID, DoorlockID"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/customer/{cccd}/scheduler [post]
func (h *CustomerHandler) AppendCustomerScheduler(c *gin.Context) {
	usu := &models.UserSchedulerUpsert{}
	cccd := c.Param("cccd")
	err := c.ShouldBind(usu)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	cus, err := h.svc.FindCustomerByCCCD(c, cccd)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get customer failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	sche := &usu.Scheduler
	sche.EmployeeID = &cus.CCCD
	sche.DoorSerialID = &usu.DoorlockID
	_, err = h.scheSvc.CreateScheduler(c, sche)

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
			UserId:     cus.CCCD,
			RfidPass:   cus.RfidPass,
			KeypadPass: cus.KeypadPass,
		},
	))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create scheduler mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.svc.AppendCustomerScheduler(c.Request.Context(), cus, usu.DoorlockID, &usu.Scheduler)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update customer failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, true)
}
