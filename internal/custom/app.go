package custom

import (
	"github.com/hinss/go-custom/framework"
	"github.com/hinss/go-custom/framework/contract"
	"github.com/hinss/go-custom/framework/provider/app"
	cfg "github.com/hinss/go-custom/framework/provider/config"
	"github.com/hinss/go-custom/framework/provider/env"
	"github.com/hinss/go-custom/framework/provider/kernel"
	"github.com/hinss/go-custom/framework/provider/log"
	"github.com/hinss/go-custom/framework/provider/log/services"
	"github.com/hinss/go-custom/internal/custom/console"
	"os"
)

// 绑定一个全局的Logger方便使用
var Logger contract.Log

func NewApp() {

	// 初始化容器
	container := framework.NewCustomContainer()
	// 获取项目根目录
	rootDir, _ := os.Getwd()
	// 绑定 App服务提供者
	container.Bind(&app.CustomAppProvider{BaseFolder: rootDir})
	// 绑定 env服务提供者
	container.Bind(&env.CustomEnvProvider{})
	// 绑定 配置服务提供者
	container.Bind(&cfg.CustomConfigProvider{})
	// 绑定 日志服务提供者
	container.Bind(&log.CustomLogServiceProvider{})
	Logger = container.MustMake(contract.LogKey).(*services.CustomRotateLog)

	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := NewHttpEngine(); err == nil {
		container.Bind(&kernel.CustomKernelProvider{HttpEngine: engine})
	}

	// 通过命令来启动
	console.RunCommand(container)
}
