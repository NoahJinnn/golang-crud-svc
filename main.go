package main

import (
	"fmt"

	"github.com/ecoprohcm/DMS_BackendServer/docs"
	"github.com/ecoprohcm/DMS_BackendServer/handlers"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/gin-gonic/gin"                 // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles"
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

func setupRouter(
	gwHandler *handlers.GatewayHandler,
	areaHandler *handlers.AreaHandler,
	dlHandler *handlers.DoorlockHandler,
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
		v1R.DELETE("/doorlock", dlHandler.DeleteDoorlock)
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
	migrateTable(db, &models.Doorlock{})
}

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
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", user, password, host, port, database)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	gwSvc := models.NewGatewaySvc(db)
	areaSvc := models.NewAreaSvc(db)
	dlSvc := models.NewDoorlockSvc(db)

	gwHdlr := handlers.NewGatewayHandler(gwSvc)
	areaHdlr := handlers.NewAreaHandler(areaSvc)
	dlHdlr := handlers.NewDoorlockHandler(dlSvc)

	if err != nil {
		panic("failed to connect database")
	}
	migrate(db)

	r := setupRouter(gwHdlr, areaHdlr, dlHdlr)
	initSwagger(r)

	r.Run(":8080")
}
