package handlers

import (
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
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

// Find student info by mssv
// @Summary Find Student By MSSV
// @Schemes
// @Description find student info by student mssv
// @Produce json
// @Param        mssv	path	string	true	"Student MSSV"
// @Success 200 {object} models.Student
// @Failure 400 {object} utils.ErrorResponse
// @Router /v1/student/{mssv} [get]
func (h *StudentHandler) FindStudentByMSSV(c *gin.Context) {
	mssv := c.Param("mssv")

	s, err := h.svc.FindStudentByMSSV(c, mssv)
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
// @Description Update student, must have "mssv" field
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
// @Description Delete student using "mssv" field
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

func (h *StudentHandler) AppendStudentScheduler(c *gin.Context) {
	usu := &models.UserSchedulerUpsert{}
	sId := c.Param("mssv")
	err := c.ShouldBind(usu)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Invalid req body",
			ErrorMsg:   err.Error(),
		})
		return
	}

	s, err := h.svc.FindStudentByMSSV(c, sId)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Get student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_SCHEDULER_C, 1, false, mqttSvc.ServerCreateRegisterPayload(*usu, s.RfidPass, s.KeypadPass, s.MSSV))
	if err := mqttSvc.HandleMqttErr(&t); err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Create scheduler mqtt failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	_, err = h.svc.AppendStudentScheduler(c.Request.Context(), s, usu.DoorlockID, &usu.Scheduler)
	if err != nil {
		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Msg:        "Update student failed",
			ErrorMsg:   err.Error(),
		})
		return
	}

	utils.ResponseJson(c, http.StatusOK, true)
}

// func (h *StudentHandler) UpdateStudentScheduler(c *gin.Context) {
// 	usu := &models.UserSchedulerUpsert{}
// 	sId := c.Param("mssv")
// 	err := c.ShouldBind(usu)
// 	if err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Invalid req body",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	s, err := h.svc.FindStudentByMSSV(c, sId)
// 	if err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Get student failed",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	t := h.mqtt.Publish(mqttSvc.TOPIC_SV_SCHEDULER_C, 1, false, mqttSvc.ServerCreateRegisterPayload(*usu, s.RfidPass, s.KeypadPass, s.MSSV))
// 	if err := mqttSvc.HandleMqttErr(&t); err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Update scheduler mqtt failed",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	_, err = h.svc.UpdateStudentScheduler(c.Request.Context(), s, usu.DoorlockID, &usu.Scheduler)
// 	if err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Update student failed",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	utils.ResponseJson(c, http.StatusOK, true)
// }

// func (h *StudentHandler) DeleteStudentScheduler(c *gin.Context) {
// 	usu := &models.UserSchedulerUpsert{}
// 	sId := c.Param("mssv")
// 	err := c.ShouldBind(usu)
// 	if err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Invalid req body",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	s, err := h.svc.FindStudentByMSSV(c, sId)
// 	if err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Get student failed",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	// t := h.mqtt.Publish(mqttSvc.TOPIC_SV_SCHEDULER_C, 1, false, mqttSvc.ServerDeleteRegisterPayload(*usu, s.MSSV))
// 	// if err := mqttSvc.HandleMqttErr(&t); err != nil {
// 	// 	utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 	// 		StatusCode: http.StatusBadRequest,
// 	// 		Msg:        "Delete scheduler mqtt failed",
// 	// 		ErrorMsg:   err.Error(),
// 	// 	})
// 	// 	return
// 	// }

// 	_, err = h.svc.DeleteStudentScheduler(c.Request.Context(), s, usu.DoorlockID, &usu.Scheduler)
// 	if err != nil {
// 		utils.ResponseJson(c, http.StatusBadRequest, &utils.ErrorResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Msg:        "Update student failed",
// 			ErrorMsg:   err.Error(),
// 		})
// 		return
// 	}

// 	utils.ResponseJson(c, http.StatusOK, true)
// }
