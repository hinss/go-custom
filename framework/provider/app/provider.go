package app

import (
	"github.com/hinss/go-custom/framework"
	"github.com/hinss/go-custom/framework/contract"
)

type CustomAppProvider struct {
	BaseFolder string
}

func (c *CustomAppProvider) Register(container framework.Container) framework.NewInstance {
	return NewCustomApp
}

func (c *CustomAppProvider) Boot(container framework.Container) error {
	return nil
}

func (c *CustomAppProvider) IsDefer() bool {
	return false
}

func (c *CustomAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, c.BaseFolder}
}

func (c *CustomAppProvider) Name() string {
	return contract.AppKey
}

