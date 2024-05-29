package controller

import (
	"fmt"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/irisx/controllerx"
	webapp "github.com/abmpio/webserver/app"
	"github.com/abmpio/xapikey"
)

type apiKeyController struct {
	controllerx.EntityController[xapikey.Aksk]
}

func newApiKeyController() *apiKeyController {
	c := &apiKeyController{}
	c.EntityController = controllerx.EntityController[xapikey.Aksk]{
		EntityService: getServiceGroup().apikeyService,
	}
	return c
}

func (c *apiKeyController) RegistRouter(webapp *webapp.Application, routerPath string) {
	c.EntityController.RouterPath = routerPath
	log.Logger.Debug(fmt.Sprintf("正在构建路由,%s...", routerPath))

	c.EntityController.RegistRouter(webapp,
		controllerx.BaseEntityControllerWithAllEndpointDisabled(true),
	)
}
