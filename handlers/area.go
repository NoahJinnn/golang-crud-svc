package handlers

import (
	"net/http"

	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type AreaHandler struct {
	svc *models.AreaSvc
}

func NewAreaHandler(svc *models.AreaSvc) *AreaHandler {
	return &AreaHandler{
		svc: svc,
	}
}

func (h *AreaHandler) CreateArea(c *gin.Context) {
	a := &models.Area{}
	err := c.ShouldBind(a)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}
	a, err = h.svc.CreateArea(a, c.Request.Context())
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create area failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, a)

}
