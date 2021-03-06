package handlers

import (
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/trancongduynguyen1997/golang-crud-svc/models"
	"github.com/trancongduynguyen1997/golang-crud-svc/mqttSvc"
	"github.com/trancongduynguyen1997/golang-crud-svc/utils"
)

type GatewayHandler struct {
	svc  *models.GatewaySvc
	mqtt mqtt.Client
}

func NewGatewayHandler(svc *models.GatewaySvc, mqtt mqtt.Client) *GatewayHandler {
	return &GatewayHandler{
		svc:  svc,
		mqtt: mqtt,
	}
}

// Find all gateways and doorlocks info
// @Summary Find All Gateway
// @Schemes
// @Description find all gateways info
// @Produce json
// @Success 200 {array} []models.Gateway
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateways [get]
func (h *GatewayHandler) FindAllGateway(c *gin.Context) {
	gwList, err := h.svc.FindAllGateway(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all gateways failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gwList)
}

// Find gateway and doorlock info by id
// @Summary Find Gateway By ID
// @Schemes
// @Description find gateway and doorlock info by gateway id
// @Produce json
// @Param        id	path	string	true	"Gateway ID"
// @Success 200 {object} models.Gateway
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway/{id} [get]
func (h *GatewayHandler) FindGatewayByID(c *gin.Context) {
	id := c.Param("id")

	gw, err := h.svc.FindGatewayByID(c, id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gw)
}

// Create gateway
// @Summary Create Gateway
// @Schemes
// @Description Create gateway. Send created info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagCreateGateway	true	"Fields need to create a gateway"
// @Success 200 {object} models.Gateway
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway [post]
func (h *GatewayHandler) CreateGateway(c *gin.Context) {
	gw := &models.Gateway{}
	err := c.ShouldBind(gw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}
	if len(gw.GatewayID) <= 0 || len(gw.Name) <= 0 {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Please fulfill these fields: name, gateway id",
			ErrorMsg:   "Missing on required fields",
		})
		return
	}

	gw, err = h.svc.CreateGateway(c.Request.Context(), gw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gw)
}

// Update gateway
// @Summary Update Gateway By ID
// @Schemes
// @Description Update gateway, must have "id" field. Send updated info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagUpateGateway	true	"Fields need to update a gateway"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway [patch]
func (h *GatewayHandler) UpdateGateway(c *gin.Context) {
	gw := &models.Gateway{}
	err := c.ShouldBind(gw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.svc.UpdateGateway(c.Request.Context(), gw)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_GATEWAY_U, 1, false, mqttSvc.ServerUpdateGatewayPayload(gw))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update gateway mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Delete gateway
// @Summary Delete Gateway By ID
// @Schemes
// @Description Delete gateway using "id" field. Send deleted info to MQTT broker
// @Accept  json
// @Produce json
// @Param	data	body	object{id=int}	true	"Gateway ID"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gateway [delete]
func (h *GatewayHandler) DeleteGateway(c *gin.Context) {
	dgw := &models.DeleteGateway{}
	err := c.ShouldBind(dgw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.svc.DeleteGateway(c.Request.Context(), dgw.GatewayID)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_GATEWAY_D, 1, false, mqttSvc.ServerDeleteGatewayPayload(dgw.GatewayID))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete gateway mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

func (h *GatewayHandler) DeleteGatewayDoorlock(c *gin.Context) {
	d := &models.Doorlock{}
	gwId := c.Param("id")
	err := c.ShouldBind(d)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	gw, err := h.svc.FindGatewayByID(c, gwId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.svc.DeleteGatewayDoorlock(c.Request.Context(), gw, d)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete gateway door failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, true)
}
