package service

import (
	"fmt"

	"github.com/abmpio/abmp/pkg/log"
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
	redisKey := s.getAkRedisKey(app, accessKey)
	redisValue := getServiceGroup().redisService.StringGet(redisKey)
	if redisValue.Err() == nil && redisValue.Exist() {
		aksk := &xapikey.Aksk{}
		err := redisValue.ToValue(aksk)
		if err == nil {
			return aksk, nil
		}
	}
	aksk, err := s._findByAk(app, accessKey)
	if err != nil {
		return nil, err
	}
	if aksk != nil {
		err = getServiceGroup().redisService.StringSet(redisKey, aksk)
		if err != nil {
			log.Logger.Warn(fmt.Sprintf("将xapikey保存到redis时出现异常,key:%s,err:%s",
				redisKey,
				err.Error()))
			return aksk, nil
		}
	}
	return aksk, nil
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
