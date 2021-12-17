package handlers

import (
	"net/http"

	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type GatewayHandler struct {
	svc *models.GatewaySvc
}

func NewGatewayHandler(svc *models.GatewaySvc) *GatewayHandler {
	return &GatewayHandler{
		svc: svc,
	}
}

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
	gw, err = h.svc.UpdateGateway(c.Request.Context(), gw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, gw)
}

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

	_, err = h.svc.DeleteGateway(c.Request.Context(), dgw)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete gateway failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, "Delete successfully")

}
