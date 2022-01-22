// Package main contains: http app, mqtt app, swagger doc, ORM setup
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ecoprohcm/DMS_BackendServer/docs"
	"github.com/ecoprohcm/DMS_BackendServer/handlers"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/ecoprohcm/DMS_BackendServer/mqttSvc"
	"github.com/gin-gonic/gin" // swagger embed files
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", user, password, host, port, database)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	migrate(db)

	// DI process
	gwSvc := models.NewGatewaySvc(db)
	areaSvc := models.NewAreaSvc(db)
	dlSvc := models.NewDoorlockSvc(db)
	glSvc := models.NewLogSvc(db)
	sSvc := models.NewStudentSvc(db)
	eSvc := models.NewEmployeeSvc(db)
	scheSvc := models.NewSchedulerSvc(db)
	cusSvc := models.NewCustomerSvc(db)

	mqttHost := os.Getenv("MQTT_HOST")
	mqttPort := os.Getenv("MQTT_PORT")

	mqttClient := mqttSvc.MqttClient(mqttHost, mqttPort, glSvc, dlSvc, gwSvc, eSvc)

	gwHdlr := handlers.NewGatewayHandler(gwSvc, mqttClient)
	areaHdlr := handlers.NewAreaHandler(areaSvc)
	dlHdlr := handlers.NewDoorlockHandler(dlSvc, mqttClient)
	glHdlr := handlers.NewGatewayLogHandler(glSvc)
	sHdlr := handlers.NewStudentHandler(sSvc, scheSvc, mqttClient)
	eHdlr := handlers.NewEmployeeHandler(eSvc, scheSvc, mqttClient)
	scheHdlr := handlers.NewSchedulerHandler(scheSvc, mqttClient)
	cusHdlr := handlers.NewCustomerHandler(cusSvc, scheSvc, mqttClient)

	// HTTP Serve
	r := setupRouter(gwHdlr, areaHdlr, dlHdlr, glHdlr, sHdlr, eHdlr, cusHdlr, scheHdlr)
	initSwagger(r)
	r.Run(":8080")
	mqttClient.Disconnect(250)
}

func initSwagger(r *gin.Engine) {
	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))
	// Programmatically set swagger info
	docs.SwaggerInfo.Title = "DMS Backend API"
	docs.SwaggerInfo.Description = "This is DMS backend server"
	docs.SwaggerInfo.Host = "http://iot.hcmue.space:8002"
	docs.SwaggerInfo.BasePath = "/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
