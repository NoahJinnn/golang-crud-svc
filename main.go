package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ecoprohcm/DMS_BackendServer/docs"
	"github.com/ecoprohcm/DMS_BackendServer/handlers"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/gin-gonic/gin" // swagger embed files
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func initSwagger(r *gin.Engine) {
	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))
	// programmatically set swagger info
	// Note: Use config later
	docs.SwaggerInfo.Title = "DMS Backend API"
	// docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Description = "This is DMS backend server"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

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

	gwHdlr := handlers.NewGatewayHandler(gwSvc)
	areaHdlr := handlers.NewAreaHandler(areaSvc)
	dlHdlr := handlers.NewDoorlockHandler(dlSvc)
	glHdlr := handlers.NewGatewayLogHandler(glSvc)

	go initMqttClient()

	// HTTP Serve
	r := setupRouter(gwHdlr, areaHdlr, dlHdlr, glHdlr)
	initSwagger(r)
	r.Run(":8080")
}
