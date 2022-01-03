package handlers

import (
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type PasswordHandler struct {
	svc  *models.PasswordSvc
	mqtt mqtt.Client
}

func NewPasswordHandler(svc *models.PasswordSvc, mqtt mqtt.Client) *PasswordHandler {
	return &PasswordHandler{
		svc:  svc,
		mqtt: mqtt,
	}
}

// Find all passwords info
// @Summary Find All Password
// @Schemes
// @Description find all passwords info
// @Produce json
// @Success 200 {array} []models.Password
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/passwords/{userId} [get]
func (h *PasswordHandler) FindAllPasswordByUserID(c *gin.Context) {
	userId := c.Param("userId")
	gwList, err := h.svc.FindAllPasswordByUserID(c, userId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all passwords failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gwList)
}

// Create password
// @Summary Create Password
// @Schemes
// @Description Create password
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagCreatePassword	true	"Fields need to create a password"
// @Success 200 {object} models.Password
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/password [post]
func (h *PasswordHandler) CreatePassword(c *gin.Context) {
	pw := &models.Password{}
	err := c.ShouldBind(pw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_PASSWORD_D, 1, false, mqttSvc.ServerCreatePasswordPayload(pw))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create password mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	pw, err = h.svc.CreatePassword(c.Request.Context(), pw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create password failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, pw)
}

// Update password
// @Summary Update Password By ID
// @Schemes
// @Description Update password, must have "id" field
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagUpdatePassword	true	"Fields need to update a password"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/password [patch]
func (h *PasswordHandler) UpdatePassword(c *gin.Context) {
	pw := &models.Password{}
	err := c.ShouldBind(pw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_PASSWORD_U, 1, false, mqttSvc.ServerUpdatePasswordPayload(pw))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update password mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.svc.UpdatePassword(c.Request.Context(), pw)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update password failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Delete password
// @Summary Delete Password By ID
// @Schemes
// @Description Delete password using "id" field
// @Accept  json
// @Produce json
// @Param	data	body	object{id=int}	true	"Password ID"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/password [delete]
func (h *PasswordHandler) DeletePassword(c *gin.Context) {
	pw := &models.Password{}
	err := c.ShouldBind(pw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_PASSWORD_D, 1, false, mqttSvc.ServerDeletePasswordPayload(pw))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete password mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.svc.DeletePassword(c.Request.Context(), pw)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete password failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)

}
