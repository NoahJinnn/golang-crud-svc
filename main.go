package main

import (
	"fmt"

	"github.com/ecoprohcm/DMS_BackendServer/handlers"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	host     = "localhost"
	port     = "1433"
	user     = "sa"
	password = "Iot@@123"
	database = "DevDB"
)

func setupRouter(handler handlers.HttpHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/", handler.CreateGateway)
	return r
}

func main() {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", user, password, host, port, database)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	handler := handlers.NewHttpHandler(db)
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Gateway{})
	r := setupRouter(*handler)
	r.Run(":8080")
}
