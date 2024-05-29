package controller

import (
	"fmt"
	"strings"

	"github.com/abmpio/irisx/controllerx"
	"github.com/abmpio/mongodbr"
	"github.com/abmpio/xapikey/options"

	"github.com/abmpio/abmp/pkg/log"
	webapp "github.com/abmpio/webserver/app"
)

func registControllers(webApp *webapp.Application) {
	newApiKeyController().RegistRouter(webApp, appendRouterPrefixPath("/xapikey"))
}

func RegistModelController[T mongodbr.IEntity](webApp *webapp.Application, routerPath string) {
	entityController := &controllerx.EntityController[T]{
		RouterPath: routerPath,
	}
	log.Logger.Debug(fmt.Sprintf("正在构建路由,%s...", entityController.RouterPath))
	entityController.RegistRouter(webApp)
}

func appendRouterPrefixPath(path string) string {
	o := options.GetOptions()
	prefixPath := o.RouterPrefixPath
	prefixPath = strings.TrimPrefix(prefixPath, "/")
	prefixPath = strings.TrimSuffix(prefixPath, "/")
	if len(prefixPath) <= 0 {
		return path
	}
	path = strings.TrimPrefix(path, "/")
	if len(path) <= 0 {
		return fmt.Sprintf("/%s", o.RouterPrefixPath)
	} else {
		return fmt.Sprintf("/%s/%s", o.RouterPrefixPath, path)
	}
}
