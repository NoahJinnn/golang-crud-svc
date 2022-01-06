package main

import (
	"github.com/ecoprohcm/DMS_BackendServer/handlers"
	"github.com/gin-gonic/gin"
)

func setupRouter(
	gwHandler *handlers.GatewayHandler,
	areaHandler *handlers.AreaHandler,
	dlHandler *handlers.DoorlockHandler,
	glHandler *handlers.GatewayLogHandler,
	// pwHdlr *handlers.PasswordHandler,
	sHdlr *handlers.StudentHandler,
	eHdlr *handlers.EmployeeHandler,
	scheHdlr *handlers.SchedulerHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())
	v1R := r.Group("/v1")
	{
		// Gateway routes
		v1R.GET("/gateways", gwHandler.FindAllGateway)
		v1R.GET("/gateway/:id", gwHandler.FindGatewayByID)
		v1R.POST("/gateway", gwHandler.CreateGateway)
		v1R.PATCH("/gateway", gwHandler.UpdateGateway)
		v1R.DELETE("/gateway", gwHandler.DeleteGateway)

		// Area routes
		v1R.GET("/areas", areaHandler.FindAllArea)
		v1R.GET("/area/:id", areaHandler.FindAreaByID)
		v1R.POST("/area", areaHandler.CreateArea)
		v1R.PATCH("/area", areaHandler.UpdateArea)
		v1R.DELETE("/area", areaHandler.DeleteArea)

		// Doorlock routes
		v1R.GET("/doorlocks", dlHandler.FindAllDoorlock)
		v1R.GET("/doorlock/:id", dlHandler.FindDoorlockByID)
		v1R.POST("/doorlock", dlHandler.CreateDoorlock)
		v1R.PATCH("/doorlock", dlHandler.UpdateDoorlock)
		v1R.PATCH("/doorlock/cmd", dlHandler.UpdateDoorlockCmd)
		v1R.DELETE("/doorlock", dlHandler.DeleteDoorlock)

		// Password routes
		// v1R.GET("/passwords/:userId", pwHdlr.FindAllPasswordByUserID)
		// v1R.POST("/password", pwHdlr.CreatePassword)
		// v1R.PATCH("/password", pwHdlr.UpdatePassword)
		// v1R.DELETE("/password", pwHdlr.DeletePassword)

		// Student routes
		v1R.GET("/students", sHdlr.FindAllStudent)
		v1R.GET("/student/:id", sHdlr.FindStudentByID)
		v1R.POST("/student", sHdlr.CreateStudent)
		v1R.PATCH("/student", sHdlr.UpdateStudent)
		v1R.DELETE("/student", sHdlr.DeleteStudent)
		v1R.POST("/student/:id/scheduler", sHdlr.AppendStudentScheduler)
		v1R.PATCH("/student/:id/scheduler", sHdlr.UpdateStudentScheduler)
		v1R.DELETE("/student/:id/scheduler", sHdlr.DeleteStudentScheduler)

		// Employee routes
		v1R.GET("/employees", eHdlr.FindAllEmployee)
		v1R.GET("/employee/:id", eHdlr.FindEmployeeByID)
		v1R.POST("/employee", eHdlr.CreateEmployee)
		v1R.PATCH("/employee", eHdlr.UpdateEmployee)
		v1R.DELETE("/employee", eHdlr.DeleteEmployee)

		// Scheduler routes
		v1R.GET("/schedulers", scheHdlr.FindAllScheduler)
		v1R.GET("/scheduler/:id", scheHdlr.FindSchedulerByID)
		v1R.POST("/scheduler", scheHdlr.CreateScheduler)
		v1R.PATCH("/scheduler", scheHdlr.UpdateScheduler)
		v1R.DELETE("/scheduler", scheHdlr.DeleteScheduler)

		// Gateway log routes
		v1R.GET("/gatewayLogs", glHandler.FindAllGatewayLog)
		v1R.GET("/gatewayLog/:id", glHandler.FindGatewayLogByID)
	}
	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Accept, Origin, Cache-Control, X-Requested-With, User-Agent, Accept-Language, Accept-Encoding")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
