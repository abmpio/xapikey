package starter

import (
	"fmt"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/app"
	"github.com/abmpio/app/cli"
	"github.com/abmpio/xapikey/middleware"
	_ "github.com/abmpio/xapikey/mongodb"
	"github.com/abmpio/xapikey/options"

	_ "github.com/abmpio/xapikey/service"

	webapp "github.com/abmpio/webserver/app"
	_ "github.com/abmpio/xapikey/controller"
)

func init() {
	fmt.Println("plugins.xapikey.starter init")

	cli.ConfigureService(serviceConfigurator)
}

func serviceConfigurator(cliApp cli.CliApplication) {
	if app.HostApplication.SystemConfig().App.IsRunInCli {
		return
	}
	webApp := app.Context.GetInstance(&webapp.Application{}).(*webapp.Application)
	if webApp != nil && !options.GetOptions().DisableXApiKey {
		middleware.UseXApiKey(webApp.APIBuilder)
		log.Logger.Info("已启用xapikey中间件")
	}
}
