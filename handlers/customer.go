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
	svc  *models.CustomerSvc
	mqtt mqtt.Client
}

func NewCustomerHandler(svc *models.CustomerSvc, mqtt mqtt.Client) *CustomerHandler {
	return &CustomerHandler{
		svc:  svc,
		mqtt: mqtt,
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

// Find customer info by id
// @Summary Find Customer By ID
// @Schemes
// @Description find customer info by customer id
// @Produce json
// @Param        id	path	string	true	"Customer ID"
// @Success 200 {object} models.Customer
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/customer/{id} [get]
func (h *CustomerHandler) FindCustomerByID(c *gin.Context) {
	id := c.Param("id")

	cus, err := h.svc.FindCustomerByID(c, id)
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
// @Description Update customer, must have "id" field
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagUpdateCustomer	true	"Fields need to update a customer"
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
	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Delete customer
// @Summary Delete Customer By ID
// @Schemes
// @Description Delete customer using "id" field
// @Accept  json
// @Produce json
// @Param	data	body	object{id=int}	true	"Customer ID"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/customer [delete]
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	dId := &models.DeleteID{}
	err := c.ShouldBind(dId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.svc.DeleteCustomer(c.Request.Context(), dId.ID)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete customer failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)

}

func (h *CustomerHandler) AppendCustomerScheduler(c *gin.Context) {
	usu := &models.UserSchedulerUpsert{}
	cId := c.Param("id")
	err := c.ShouldBind(usu)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	cus, err := h.svc.FindCustomerByID(c, cId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get customer failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_SCHEDULER_C, 1, false, mqttSvc.ServerUpsertRegisterPayload(*usu, cus.RfidPass, cus.KeypadPass, cus.CCCD))
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

func (h *CustomerHandler) UpdateCustomerScheduler(c *gin.Context) {
	usu := &models.UserSchedulerUpsert{}
	cId := c.Param("id")
	err := c.ShouldBind(usu)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	cus, err := h.svc.FindCustomerByID(c, cId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get customer failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_SCHEDULER_C, 1, false, mqttSvc.ServerUpsertRegisterPayload(*usu, cus.RfidPass, cus.KeypadPass, cus.CCCD))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update scheduler mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.svc.UpdateCustomerScheduler(c.Request.Context(), cus, usu.DoorlockID, &usu.Scheduler)
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

func (h *CustomerHandler) DeleteCustomerScheduler(c *gin.Context) {
	usu := &models.UserSchedulerUpsert{}
	cId := c.Param("id")
	err := c.ShouldBind(usu)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	cus, err := h.svc.FindCustomerByID(c, cId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get customer failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_SCHEDULER_C, 1, false, mqttSvc.ServerDeleteRegisterPayload(*usu, cus.CCCD))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete scheduler mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.svc.DeleteCustomerScheduler(c.Request.Context(), cus, usu.DoorlockID, &usu.Scheduler)
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
