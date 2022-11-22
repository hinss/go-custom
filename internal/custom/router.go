package custom

import (
	"github.com/hinss/go-custom/framework/gin"
	"github.com/hinss/go-custom/framework/middleware"
	ginSwagger "github.com/hinss/go-custom/framework/middleware/gin-swagger"
	"github.com/hinss/go-custom/framework/middleware/gin-swagger/swaggerFiles"
	"net/http"
	"os"
	"path"
)

func registerRouter(e *gin.Engine) {

	// 文件服务路由
	rootWd, _ := os.Getwd()
	rootPath := path.Join(rootWd, "static")
	e.Use(middleware.ServeRoot("/", rootPath))

	//e.StaticFS("/static", http.Dir(config.Static))
	e.GET("/health", healthCheckHandler)

	// swagger 中间件
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}


func healthCheckHandler(ctx *gin.Context) {

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.WriteString("health")
}
