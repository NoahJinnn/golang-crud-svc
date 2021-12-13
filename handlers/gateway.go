package handlers

import (
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/gin-gonic/gin"
)

func (h *HttpHandler) CreateGateway(c *gin.Context) {
	gw := models.Gateway{
		GatewayID: "test",
		AreaID:    1,
		Name:      "Test GW",
	}
	gw.CreateGateway(h.DB)

}
