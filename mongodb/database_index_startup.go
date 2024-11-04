package mongodb

import (
	"fmt"
	"strings"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/app"
	"github.com/abmpio/app/cli"
	"github.com/abmpio/entity/mongodb"
	"github.com/abmpio/mongodbr"
	"github.com/abmpio/xapikey"
	"github.com/abmpio/xapikey/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func mustCreateIndexConfigurator(cliApp cli.CliApplication) {
	if app.HostApplication.SystemConfig().App.IsRunInCli {
		return
	}
	mustCreateIndexes()
}

func mustCreateIndexes() {
	//init indexes
	mustCreateIndexForEntity[xapikey.Aksk]([]mongo.IndexModel{
		{Keys: bson.M{"creationTime": -1}},
		{Keys: bson.D{
			{Key: "tenantId", Value: 1},
			{Key: "app", Value: 1},
			{Key: "accessKey", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "tenantId", Value: 1},
			{Key: "app", Value: 1},
			{Key: "alias", Value: 1},
		}},
	})
}

func mustCreateIndexForEntity[T mongodbr.IEntity](indexes []mongo.IndexModel) {
	collectionName := mongodb.GetCollectionName(new(T))
	database := getDatabase()
	col := mongodbr.NewMongoCol(database.Collection(collectionName))
	names, err := col.CreateIndexes(indexes)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("在为表 %s创建索引时出现错误,%s", collectionName, err.Error()))
	} else {
		log.Logger.Debug(fmt.Sprintf("表 %s索引创建成功,%s", collectionName, strings.Join(names, ",")))
	}
}

func getDatabase() *mongo.Database {
	o := options.GetOptions()
	return mongodbr.GetDatabaseByKey(o.MongodbClientKey, o.DatabaseName)
}
