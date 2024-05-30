package starter

import (
	"github.com/abmpio/app"
	_ "github.com/abmpio/xapikey/mongodb"

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
	})
}
