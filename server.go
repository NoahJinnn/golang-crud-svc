package main

import (
	"github.com/gin-gonic/gin"
	"github.com/trancongduynguyen1997/golang-crud-svc/handlers"
)

func setupRouter(
	gwHandler *handlers.GatewayHandler,
	areaHandler *handlers.AreaHandler,
	dlHandler *handlers.DoorlockHandler,
	glHandler *handlers.GatewayLogHandler,
	sHdlr *handlers.StudentHandler,
	eHdlr *handlers.EmployeeHandler,
	cusHdlr *handlers.CustomerHandler,
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
		v1R.DELETE("/gateway/:id/doorlock", gwHandler.DeleteGatewayDoorlock)

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

		// Student routes
		v1R.GET("/students", sHdlr.FindAllStudent)
		v1R.GET("/student/:mssv", sHdlr.FindStudentByMSSV)
		v1R.POST("/student", sHdlr.CreateStudent)
		v1R.PATCH("/student", sHdlr.UpdateStudent)
		v1R.DELETE("/student", sHdlr.DeleteStudent)
		v1R.POST("/student/:mssv/scheduler", sHdlr.AppendStudentScheduler)

		// Employee routes
		v1R.GET("/employees", eHdlr.FindAllEmployee)
		v1R.GET("/employee/:msnv", eHdlr.FindEmployeeByMSNV)
		v1R.POST("/employee", eHdlr.CreateEmployee)
		v1R.PATCH("/employee", eHdlr.UpdateEmployee)
		v1R.DELETE("/employee", eHdlr.DeleteEmployee)
		v1R.POST("/employee/:msnv/scheduler", eHdlr.AppendEmployeeScheduler)

		// Customer routes
		v1R.GET("/customers", cusHdlr.FindAllCustomer)
		v1R.GET("/customer/:cccd", cusHdlr.FindCustomerByCCCD)
		v1R.POST("/customer", cusHdlr.CreateCustomer)
		v1R.PATCH("/customer", cusHdlr.UpdateCustomer)
		v1R.DELETE("/customer", cusHdlr.DeleteCustomer)
		v1R.POST("/customer/:cccd/scheduler", cusHdlr.AppendCustomerScheduler)

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
