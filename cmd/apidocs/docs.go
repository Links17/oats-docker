package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.GET("/api-ref/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8091")
}
