package kernel

import (
	"github.com/hinss/go-custom/framework/gin"
	"net/http"
)

// 引擎服务
type CustomKernelService struct {
	engine *gin.Engine
}

// 初始化web引擎服务实例
func NewCustomKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &CustomKernelService{engine: httpEngine}, nil
}

// 返回web引擎
func (s *CustomKernelService) HttpEngine() http.Handler {
	return s.engine
}
