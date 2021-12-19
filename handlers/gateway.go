package handlers

import (
	"net/http"
	"strconv"

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

func (h *GatewayHandler) FindGatewayByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req id",
			ErrorMsg:   err.Error(),
		})
		return
	}
	gw, err := h.svc.FindGatewayByID(c, uint(id))
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
	if len(gw.MacID) <= 0 || len(gw.Name) <= 0 {
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
	g := &models.Gateway{}
	err := c.ShouldBind(g)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.svc.DeleteGateway(c.Request.Context(), g)
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
