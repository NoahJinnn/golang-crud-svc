package handlers

import (
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type DoorlockHandler struct {
	svc  *models.DoorlockSvc
	mqtt mqtt.Client
}

func NewDoorlockHandler(svc *models.DoorlockSvc, mqtt mqtt.Client) *DoorlockHandler {
	return &DoorlockHandler{
		svc:  svc,
		mqtt: mqtt,
	}
}

// Find all doorlocks info
// @Summary Find All Doorlock
// @Schemes
// @Description find all doorlocks info
// @Produce json
// @Success 200 {array} []models.Doorlock
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlocks [get]
func (h *DoorlockHandler) FindAllDoorlock(c *gin.Context) {
	dlList, err := h.svc.FindAllDoorlock(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all doorlocks failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, dlList)
}

// Find doorlock info by id
// @Summary Find Doorlock By ID
// @Schemes
// @Description find doorlock info by doorlock id
// @Produce json
// @Param        id	path	string	true	"Doorlock ID"
// @Success 200 {object} models.Doorlock
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlock/{id} [get]
func (h *DoorlockHandler) FindDoorlockByID(c *gin.Context) {
	id := c.Param("id")

	dl, err := h.svc.FindDoorlockByID(c, id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get doorlock failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, dl)
}

// Create doorlock
// @Summary Create Doorlock
// @Schemes
// @Description Create doorlock
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagCreateDoorlock	true	"Fields need to create a doorlock"
// @Success 200 {object} models.Doorlock
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlock [post]
func (h *DoorlockHandler) CreateDoorlock(c *gin.Context) {
	dl := &models.Doorlock{}
	err := c.ShouldBind(dl)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}
	if len(dl.Location) <= 0 {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Please fulfill these fields: location",
			ErrorMsg:   "Missing on required fields",
		})
		return
	}

	t := h.mqtt.Publish(string(mqttSvc.TOPIC_SV_DOORLOCK_C), 1, false,
		mqttSvc.ServerCreateDoorlockPayload(dl),
	)
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create doorlock mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	dl, err = h.svc.CreateDoorlock(c.Request.Context(), dl)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create doorlock failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, dl)

}

// Update doorlock
// @Summary Update Doorlock By ID
// @Schemes
// @Description Update doorlock, must have "id" field
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagUpdateDoorlock	true	"Fields need to update a doorlock"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlock [patch]
func (h *DoorlockHandler) UpdateDoorlock(c *gin.Context) {
	dl := &models.Doorlock{}
	err := c.ShouldBind(dl)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_DOORLOCK_U, 1, false, mqttSvc.ServerUpdateDoorlockPayload(dl))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update doorlock mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.svc.UpdateDoorlock(c.Request.Context(), dl)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update doorlock failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Update doorlock state
// @Summary Update Doorlock State By ID
// @Schemes
// @Description Update doorlock state, must have "id" field
// @Accept  json
// @Produce json
// @Param	data	body	models.DoorlockCmd	true	"Fields need to update a doorlock state"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlock/state [patch]
func (h *DoorlockHandler) UpdateDoorlockCmd(c *gin.Context) {
	dl := &models.DoorlockCmd{}
	err := c.ShouldBind(dl)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(string(mqttSvc.TOPIC_SV_DOORLOCK_CMD), 1, false,
		mqttSvc.ServerCmdDoorlockPayload(dl.GatewayID, dl.DoorSerialID, dl.State))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Execute doorlock command failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	// TODO: Guarantee mqtt req/res
	// isMqttReps := waitForMqttDoorlockResponse(c, 60)
	// if !isMqttReps {
	// 	utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
	// 		StatusCode: http.StatusBadRequest,
	// 		Msg:        "Mqtt response is too long",
	// 		ErrorMsg:   err.Error(),
	// 	})
	// 	return
	// }

	isSuccess, err := h.svc.UpdateDoorlockState(c.Request.Context(), dl)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update doorlock failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Delete doorlock
// @Summary Delete Doorlock By ID
// @Schemes
// @Description Delete doorlock using "id" field
// @Accept  json
// @Produce json
// @Param	data	body	models.DoorlockDelete	true	"Doorlock Delete payload"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/doorlock [delete]
func (h *DoorlockHandler) DeleteDoorlock(c *gin.Context) {
	dl := &models.DoorlockDelete{}
	err := c.ShouldBind(dl)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_DOORLOCK_D, 1, false,
		mqttSvc.ServerDeleteDoorlockPayload(dl))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete doorlock mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	// TODO: Guarantee mqtt req/res
	// isMqttReps := waitForMqttDoorlockResponse(c, 60)
	// if !isMqttReps {
	// 	utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
	// 		StatusCode: http.StatusBadRequest,
	// 		Msg:        "Mqtt response is too long",
	// 		ErrorMsg:   err.Error(),
	// 	})
	// 	return
	// }

	isSuccess, err := h.svc.DeleteDoorlock(c.Request.Context(), dl.DoorSerialID)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete doorlock failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)

}

// TODO: Guarantee mqtt req/res
// func waitForMqttDoorlockResponse(c *gin.Context, waitSecond uint) bool {
// 	// Clear previous check read state
// 	select {
// 	case <-mqttSvc.DoorlockStateCheck:
// 		fmt.Println("[MQTT-INFO] Clear doorlock state check flag")
// 	default:
// 	}

// 	for start := time.Now(); time.Since(start) < time.Duration(waitSecond)*time.Second; {
// 		select {
// 		case <-mqttSvc.DoorlockStateCheck:
// 			return true
// 		default:
// 			fmt.Println("[MQTT-INFO] Checking doorlock state from mqtt")

// 		}
// 	}
// 	return false
// }
