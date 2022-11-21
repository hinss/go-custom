package env

import (
	"github.com/hinss/go-custom/framework"
	"github.com/hinss/go-custom/framework/contract"
)

type CustomEnvProvider struct {
	Folder string
}

// Register registe a new function for make a service instance
func (provider *CustomEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewCustomEnv
}

// Boot will called when the service instantiate
func (provider *CustomEnvProvider) Boot(c framework.Container) error {
	app := c.MustMake(contract.AppKey).(contract.App)
	provider.Folder = app.BaseFolder()
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *CustomEnvProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *CustomEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.Folder}
}

/// Name define the name for this service
func (provider *CustomEnvProvider) Name() string {
	return contract.EnvKey
}
