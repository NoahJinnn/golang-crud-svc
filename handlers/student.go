package handlers

import (
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/utils"
	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	svc  *models.StudentSvc
	mqtt mqtt.Client
}

func NewStudentHandler(svc *models.StudentSvc, mqtt mqtt.Client) *StudentHandler {
	return &StudentHandler{
		svc:  svc,
		mqtt: mqtt,
	}
}

// Find all students info
// @Summary Find All Student
// @Schemes
// @Description find all students info
// @Produce json
// @Success 200 {array} []models.Student
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/students [get]
func (h *StudentHandler) FindAllStudent(c *gin.Context) {
	sList, err := h.svc.FindAllStudent(c)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get all students failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, sList)
}

// Find student info by id
// @Summary Find Student By ID
// @Schemes
// @Description find student info by student id
// @Produce json
// @Param        id	path	string	true	"Student ID"
// @Success 200 {object} models.Student
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/student/{id} [get]
func (h *StudentHandler) FindStudentByID(c *gin.Context) {
	id := c.Param("id")

	s, err := h.svc.FindStudentByID(c, id)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, s)
}

// Create student
// @Summary Create Student
// @Schemes
// @Description Create student
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagCreateStudent	true	"Fields need to create a student"
// @Success 200 {object} models.Student
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/student [post]
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	s := &models.Student{}
	err := c.ShouldBind(s)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.svc.CreateStudent(c.Request.Context(), s)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, s)
}

// Update student
// @Summary Update Student By ID
// @Schemes
// @Description Update student, must have "id" field
// @Accept  json
// @Produce json
// @Param	data	body	models.SwagUpdateStudent	true	"Fields need to update a student"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/student [patch]
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	s := &models.Student{}
	err := c.ShouldBind(s)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	isSuccess, err := h.svc.UpdateStudent(c.Request.Context(), s)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)
}

// Delete student
// @Summary Delete Student By ID
// @Schemes
// @Description Delete student using "id" field
// @Accept  json
// @Produce json
// @Param	data	body	object{id=int}	true	"Student ID"
// @Success 200 {boolean} true
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/student [delete]
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
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

	isSuccess, err := h.svc.DeleteStudent(c.Request.Context(), dId.ID)
	if err != nil || !isSuccess {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Delete student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}
	utils.ResponseJson(c, http.StatusOK, isSuccess)

}
