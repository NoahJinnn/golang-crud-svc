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

func setupRouter(gwHandler *handlers.GatewayHandler, areaHandler *handlers.AreaHandler) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())
	v1R := r.Group("/v1")
	{
		v1R.GET("/gateways", gwHandler.FindAllGateway)
		v1R.POST("/gateway", gwHandler.CreateGateway)
		v1R.PATCH("/gateway", gwHandler.UpdateGateway)
		v1R.DELETE("/gateway", gwHandler.DeleteGateway)
		v1R.POST("/area", areaHandler.CreateArea)
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

func migrateTable(db *gorm.DB, t interface{}) {
	if !db.Migrator().HasTable(t) {
		db.Migrator().CreateTable(t)
	}
}

func migrate(db *gorm.DB) {
	migrateTable(db, &models.Gateway{})
	migrateTable(db, &models.Area{})
}

func main() {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", user, password, host, port, database)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	gwSvc := models.NewGatewaySvc(db)
	areaSvc := models.NewAreaSvc(db)

	gwHdlr := handlers.NewGatewayHandler(gwSvc)
	areaHdlr := handlers.NewAreaHandler(areaSvc)

	if err != nil {
		panic("failed to connect database")
	}
	migrate(db)

	r := setupRouter(gwHdlr, areaHdlr)
	r.Run(":8080")
}
