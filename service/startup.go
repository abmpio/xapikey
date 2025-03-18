package service

import (
	"github.com/abmpio/app"
	"github.com/abmpio/app/cli"
	"github.com/abmpio/entity"
	"github.com/abmpio/entity/mongodb"
	"github.com/abmpio/xapikey"
	"github.com/abmpio/xapikey/options"
)

func init() {
	cli.ConfigureService(serviceConfigurator)
}

func serviceConfigurator(cliApp cli.CliApplication) {
	database := getDatabase()
	//注册服务到ioc容器中
	xapiKeyService := newApiKeyService(database.GetRepository(new(xapikey.Aksk)))
	app.Context.RegistInstanceAs(xapiKeyService, new(entity.IEntityService[xapikey.Aksk]))
	app.Context.RegistInstanceAs(xapiKeyService, new(xapikey.IAkskService))

}

func getDatabase() *mongodb.Database {
	options := options.GetOptions()
	return mongodb.GetDatabase(options.MongodbClientKey, options.DatabaseName)
}
