package handlers

import (
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type SchedulerHandler struct {
	svc  *models.SchedulerSvc
	mqtt mqtt.Client
}

func NewSchedulerHandler(svc *models.SchedulerSvc, mqtt mqtt.Client) *SchedulerHandler {
	return &SchedulerHandler{
		svc:  svc,
		mqtt: mqtt,
	}
}

// Find all scheduler info
// @Summary Find All Scheduler
// @Schemes
// @Description find all scheduler info
// @Produce json
// @Success 200 {array} []models.Scheduler
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/schedulers [get]
func (h *SchedulerHandler) FindAllScheduler(c *gin.Context) {
	sList, err := h.svc.FindAllScheduler(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all scheduler failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, sList)
}

// Find scheduler info by id
// @Summary Find Scheduler By ID
// @Schemes
// @Description find scheduler info by scheduler id
// @Produce json
// @Param        id	path	string	true	"Scheduler ID"
// @Success 200 {object} models.Scheduler
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/scheduler/{id} [get]
func (h *SchedulerHandler) FindSchedulerByID(c *gin.Context) {
	id := c.Param("id")

	s, err := h.svc.FindSchedulerByID(c, id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get scheduler failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, s)
}

// Create scheduler
// @Summary Create Scheduler
// @Schemes
// @Description Create scheduler
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagCreateScheduler	true	"Fields need to create a scheduler"
// @Success 200 {object} models.Scheduler
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/scheduler [post]
func (h *SchedulerHandler) CreateScheduler(c *gin.Context) {
	s := &models.Scheduler{}
	err := c.ShouldBind(s)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.svc.CreateScheduler(c.Request.Context(), s)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create scheduler failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, s)
}

// Update scheduler
// @Summary Update Scheduler By ID
// @Schemes
// @Description Update scheduler, must have "id" field
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagUpdateScheduler	true	"Fields need to update a scheduler"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/scheduler [patch]
func (h *SchedulerHandler) UpdateScheduler(c *gin.Context) {
	s := &models.Scheduler{}
	err := c.ShouldBind(s)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.svc.UpdateScheduler(c.Request.Context(), s)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update scheduler failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Delete scheduler
// @Summary Delete Scheduler By ID
// @Schemes
// @Description Delete scheduler using "id" field
// @Accept  json
// @Produce json
// @Param	data	body	object{id=int}	true	"Scheduler ID"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/scheduler [delete]
func (h *SchedulerHandler) DeleteScheduler(c *gin.Context) {
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

	isSuccess, err := h.svc.DeleteScheduler(c.Request.Context(), dId.ID)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete scheduler failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)

}
