package service

import (
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

func (s *apiKeyService) FindByAk(app string, accessKey string) ([]*xapikey.Aksk, error) {
	filter := bson.M{
		"app":       app,
		"accessKey": accessKey,
	}
	return s.FindList(filter)
}
