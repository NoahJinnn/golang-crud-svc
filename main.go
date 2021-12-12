package main

import (
	"fmt"

	"github.com/ecoprohcm/DMS_BackendServer/handlers"
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

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

func main() {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", user, password, host, port, database)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	handlers.NewHttpHandler(db)
	if err != nil {
		panic("failed to connect database")
	}

	r := setupRouter()
	r.Run(":8080")
}
