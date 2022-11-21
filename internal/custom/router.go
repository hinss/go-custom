package custom

import (
	"github.com/hinss/go-custom/framework/gin"
	"net/http"
)

func registerRouter(e *gin.Engine) {

	//e.StaticFS("/static", http.Dir(config.Static))
	e.GET("/health", healthCheckHandler)
}


func healthCheckHandler(ctx *gin.Context) {

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.WriteString("health")
}
