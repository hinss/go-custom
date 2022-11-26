package tests

import (
	"github.com/hinss/go-custom/framework"
	"github.com/hinss/go-custom/framework/provider/app"
)

const (
	BasePath = "/Users/yejianfeng/Documents/UGit/coredemo/"
)

func InitBaseContainer() framework.Container {
	// 初始化服务容器
	container := framework.NewCustomContainer()
	// 绑定App服务提供者
	container.Bind(&app.CustomAppProvider{BaseFolder: BasePath})
	// 后续初始化需要绑定的服务提供者...
	//container.Bind(&env.CustomTestingEnvProvider{})
	return container
}
