package mongodb

import (
	"github.com/abmpio/app"
	"github.com/abmpio/app/cli"
	"github.com/abmpio/entity/mongodb"
	"github.com/abmpio/xapikey"
	"github.com/abmpio/xapikey/options"
)

func repositoryConfigurator(cliApp cli.CliApplication) {
	if app.HostApplication.SystemConfig().App.IsRunInCli {
		return
	}

	o := options.GetOptions()

	mongodb.RegistEntityRepositoryOption[xapikey.Aksk](o.MongodbClientKey, o.DatabaseName)
}
