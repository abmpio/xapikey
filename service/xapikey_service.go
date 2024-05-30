package service

import (
	"fmt"

	"github.com/abmpio/entity"
	"github.com/abmpio/mongodbr"
	"github.com/abmpio/xapikey"
	"go.mongodb.org/mongo-driver/bson"
)

type apiKeyService struct {
	entity.IEntityService[xapikey.Aksk]
}

var _ xapikey.IAkskService = (*apiKeyService)(nil)

func newApiKeyService(repository mongodbr.IRepository) xapikey.IAkskService {
	s := &apiKeyService{
		IEntityService: entity.NewEntityService[xapikey.Aksk](repository),
	}
	return s
}

func (s *apiKeyService) FindByAk(app string, accessKey string) (*xapikey.Aksk, error) {
	redisValue := getServiceGroup().redisService.StringGet(s.getAkRedisKey(app, accessKey))
	aksk : *xapikey.Aksk
	var err error
	if redisValue.Err() == nil && redisValue.Exist() {
		redisValue.ToValue(adsk)
	}
	filter := bson.M{
		"app":       app,
		"accessKey": accessKey,
	}
	return s.FindOne(filter)
}

func (s *apiKeyService) _findByAk(app string, accessKey string) (*xapikey.Aksk, error) {
	filter := bson.M{
		"app":       app,
		"accessKey": accessKey,
	}
	return s.FindOne(filter)
}

func (s *apiKeyService) getAkRedisKey(app string, ak string) string {
	return fmt.Sprintf("ak:%s:%s", app, ak)
}
