package controller

import (
	"sync"

	"github.com/abmpio/app"
	"github.com/abmpio/xapikey"
)

var (
	_serviceFactory             *serviceGroup
	_serviceFactoryInstanceOnce sync.Once
)

type serviceGroup struct {
	apikeyService xapikey.IAkskService
}

// 获取ServiceGroup实例
func getServiceGroup() *serviceGroup {
	if _serviceFactory != nil {
		return _serviceFactory
	}
	_serviceFactoryInstanceOnce.Do(func() {
		_serviceFactory = newServiceGroup()
	})
	return _serviceFactory
}

// 创建一个新的实例
func newServiceGroup() *serviceGroup {
	serviceFactory := &serviceGroup{
		apikeyService: app.Context.GetInstance(new(xapikey.IAkskService)).(xapikey.IAkskService),
	}
	return serviceFactory
}
