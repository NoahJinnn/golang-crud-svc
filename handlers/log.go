package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trancongduynguyen1997/golang-crud-svc/models"
	"github.com/trancongduynguyen1997/golang-crud-svc/utils"
)

type GatewayLogHandler struct {
	svc *models.LogSvc
}

func NewGatewayLogHandler(svc *models.LogSvc) *GatewayLogHandler {
	return &GatewayLogHandler{
		svc: svc,
	}
}

// Find all gateway logs info
// @Summary Find All GatewayLog
// @Schemes
// @Description find all gateway logs info
// @Produce json
// @Success 200 {array} []models.GatewayLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gatewayLogs [get]
func (h *GatewayLogHandler) FindAllGatewayLog(c *gin.Context) {
	glList, err := h.svc.FindAllGatewayLog(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all gateway logs failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, glList)
}

// Find gateway log info by id
// @Summary Find GatewayLog By ID
// @Schemes
// @Description find gateway log info by id
// @Produce json
// @Param        id	path	string	true	"GatewayLog ID"
// @Success 200 {object} models.GatewayLog
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/gatewayLog/{id} [get]
func (h *GatewayLogHandler) FindGatewayLogByID(c *gin.Context) {
	id := c.Param("id")
	gl, err := h.svc.FindGatewayLogByID(c, id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get gateway log failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gl)
}
