package starter

import (
	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/app"
	"github.com/abmpio/xapikey/middleware"
	_ "github.com/abmpio/xapikey/mongodb"
	"github.com/abmpio/xapikey/options"

	_ "github.com/abmpio/xapikey/service"

	webapp "github.com/abmpio/webserver/app"
	_ "github.com/abmpio/xapikey/controller"
)

func init() {
	app.RegisterOneStartupAction(xapiKeyMiddlewareStartupAction)
}

func xapiKeyMiddlewareStartupAction(webApp *webapp.Application) app.IStartupAction {
	return app.NewStartupAction(func() {
		if app.HostApplication.SystemConfig().App.IsRunInCli {
			return
		}
		if !options.GetOptions().DisableXApiKey {
			middleware.UseXApiKey(webApp.APIBuilder)
			log.Logger.Info("已启用xapikey中间件")
		}
	})
}
