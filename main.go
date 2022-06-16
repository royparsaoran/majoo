package main

import (
	"majoo/conn"
	"majoo/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Swagger for Authentication API
// @version 1.0
// @description Swagger for backend API service
// @description Get the Bearer token on the Authentication Service
// @description JSON Link: <a href=/swagger/doc.json>docs.json</a>

// @securityDefinitions.apikey  BearerAuth
// @in header
// @name Authorization

// @docExpansion none

func main() {

	conn.DBConn = conn.DBEstablish()
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)

	router.GET("/merchant/:merchant_id/omzet", controller.GetMerchantOmzet)
	router.GET("/merchant/:merchant_id/outlet/:outlet_id/omzet", controller.GetOutletOmzet)

	router.Run(":8092")
}
