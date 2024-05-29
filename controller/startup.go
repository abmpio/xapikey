package controller

import (
	"github.com/abmpio/app"
	webapp "github.com/abmpio/webserver/app"
	"github.com/abmpio/xapikey/options"
)

func init() {
	app.RegisterOneStartupAction(routerStartupAction)
}

func routerStartupAction(webApp *webapp.Application) app.IStartupAction {
	return app.NewStartupAction(func() {
		if app.HostApplication.SystemConfig().App.IsRunInCli {
			return
		}
		option := options.GetOptions()
		if !option.DisableControllerRegist {
			registControllers(webApp)
		}
	})
}
