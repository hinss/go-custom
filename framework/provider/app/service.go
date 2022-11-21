package app

import (
	"flag"
	"github.com/hinss/go-custom/framework"
	"github.com/hinss/go-custom/framework/util"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"path/filepath"
)

type CustomApp struct {
	container  framework.Container // 服务容器
	baseFolder string              // 基础路径
	appId      string              // 表示当前这个app的唯一id, 可以用于分布式锁等

	configMap map[string]string // 配置加载
}

func NewCustomApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}

	// 有两个参数，一个是容器，一个是baseFolder
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)
	// 如果没有设置，则使用参数
	if baseFolder == "" {
		flag.StringVar(&baseFolder, "base_folder", "", "base_folder参数, 默认为当前路径")
		flag.Parse()
	}
	appId := uuid.New().String()
	return &CustomApp{baseFolder: baseFolder, container: container, appId: appId}, nil
}

func (c CustomApp) AppID() string {
	return c.appId
}

// Version 实现版本
func (c CustomApp) Version() string {
	return "1.0.0"
}

// BaseFolder 表示基础目录，可以代表开发场景的目录，也可以代表运行时候的目录
func (c CustomApp) BaseFolder() string {
	if c.baseFolder != "" {
		return c.baseFolder
	}
	return util.GetExecDirectory()
}

// ConfigFolder  表示配置文件地址
func (c CustomApp) ConfigFolder() string {
	if val, ok := c.configMap["config_folder"]; ok {
		return val
	}
	return filepath.Join(c.BaseFolder(), "config")
}

// LogFolder 表示日志存放地址
func (c CustomApp) LogFolder() string {
	if val, ok := c.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(c.StorageFolder(), "log")
}

func (c CustomApp) HttpFolder() string {
	if val, ok := c.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(c.BaseFolder(), "app", "http")
}

func (c CustomApp) ConsoleFolder() string {
	if val, ok := c.configMap["console_folder"]; ok {
		return val
	}
	return filepath.Join(c.BaseFolder(), "app", "console")
}

func (c CustomApp) StorageFolder() string {
	if val, ok := c.configMap["storage_folder"]; ok {
		return val
	}
	return filepath.Join(c.BaseFolder(), "storage")
}

// ProviderFolder 定义业务自己的服务提供者地址
func (c CustomApp) ProviderFolder() string {
	if val, ok := c.configMap["provider_folder"]; ok {
		return val
	}
	return filepath.Join(c.BaseFolder(), "app", "provider")
}

// MiddlewareFolder 定义业务自己定义的中间件
func (c CustomApp) MiddlewareFolder() string {
	if val, ok := c.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(c.HttpFolder(), "middleware")
}

// CommandFolder 定义业务定义的命令
func (c CustomApp) CommandFolder() string {
	if val, ok := c.configMap["command_folder"]; ok {
		return val
	}
	return filepath.Join(c.ConsoleFolder(), "command")
}

// RuntimeFolder 定义业务的运行中间态信息
func (c CustomApp) RuntimeFolder() string {
	if val, ok := c.configMap["runtime_folder"]; ok {
		return val
	}
	return filepath.Join(c.StorageFolder(), "runtime")
}

// TestFolder 定义测试需要的信息
func (c CustomApp) TestFolder() string {
	if val, ok := c.configMap["test_folder"]; ok {
		return val
	}
	return filepath.Join(c.BaseFolder(), "test")
}


func (c CustomApp) AppFolder() string {
	if val, ok := c.configMap["app"]; ok {
		return val
	}
	return filepath.Join(c.BaseFolder(), "internal")
}


// LoadAppConfig 加载配置map
func (c CustomApp) LoadAppConfig(kv map[string]string) {
	for key, val := range kv {
		c.configMap[key] = val
	}
}


