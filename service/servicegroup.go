package service

import (
	"sync"

	"github.com/abmpio/app"
	redis "github.com/abmpio/redisx"
	"github.com/abmpio/xapikey"
)

var (
	_serviceFactory             *serviceGroup
	_serviceFactoryInstanceOnce sync.Once
)

type serviceGroup struct {
	redisService  redis.IRedisService
	apiKeyService xapikey.IAkskService
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
		redisService: app.Context.GetInstance(new(redis.IRedisService)).(redis.IRedisService),

		apiKeyService: app.Context.GetInstance(new(xapikey.IAkskService)).(xapikey.IAkskService),
	}
	return serviceFactory
}

// get IAkskService instance
func ApiKeyService() xapikey.IAkskService {
	return getServiceGroup().apiKeyService
}
