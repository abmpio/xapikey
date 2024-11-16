package service

import (
	"fmt"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/entity"
	"github.com/abmpio/mongodbr"
	"github.com/abmpio/xapikey"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s *apiKeyService) FindByAk(app string, accessKey string) (aksk *xapikey.Aksk, err error) {
	key := getRedisKeyForXApiKey(app, accessKey)
	redisValue := getServiceGroup().redisService.StringGet(key)
	if redisValue.Err() != nil {
		// read redis error
		log.Logger.Warn(fmt.Sprintf("从redis中读取aksk数据时出现异常,key:%s,err:%s",
			key,
			redisValue.Err().Error()))
	} else {
		if redisValue.Exist() {
			aksk = &xapikey.Aksk{}
			err = redisValue.ToValue(aksk)
			if err == nil {
				return aksk, nil
			} else {
				log.Logger.Warn(fmt.Sprintf("将redis中读取的aksk数据反序列化成aksk对象时出现异常,key:%s,err:%s",
					key,
					err.Error()))
			}
		}
	}
	aksk, err = s._findByAk(app, accessKey)
	if err != nil {
		return nil, err
	}
	if aksk == nil {
		return aksk, nil
	}
	err = getServiceGroup().redisService.StringSet(key, aksk)
	if err != nil {
		log.Logger.Warn(fmt.Sprintf("将xapikey保存到redis时出现异常,key:%s,err:%s",
			key,
			err.Error()))
	}
	return aksk, nil
}

func (s *apiKeyService) Create(item interface{}) (*xapikey.Aksk, error) {
	aksk, ok := item.(*xapikey.Aksk)
	if !ok {
		return nil, fmt.Errorf("item必须是aksk对象")
	}
	err := aksk.Validate()
	if err != nil {
		return nil, err
	}
	list, err := s.FindList(bson.M{
		"app":       aksk.App,
		"creatorId": aksk.CreatorId,
		"alias":     aksk.Alias,
	})
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return nil, fmt.Errorf("已经存在着名称为: %s 的其它记录", aksk.Alias)
	}

	aksk.BeforeCreate()
	aksk, err = s.IEntityService.Create(aksk)
	if err != nil {
		return nil, err
	}
	return aksk, nil
}

func (s *apiKeyService) UpdateFields(id primitive.ObjectID, update map[string]interface{}) error {
	item, err := s.FindById(id)
	if err != nil {
		return err
	}
	if item == nil {
		return fmt.Errorf("无效的数据")
	}
	err = s.IEntityService.UpdateFields(id, update)
	if err != nil {
		return err
	}
	s._deleteRedisKey(item.App, item.AccessKey)
	return nil
}

func (s *apiKeyService) Delete(id primitive.ObjectID) error {
	item, err := s.FindById(id)
	if err != nil {
		return err
	}
	if item == nil {
		return fmt.Errorf("无效的数据")
	}

	err = s.IEntityService.Delete(id)
	if err != nil {
		return err
	}
	s._deleteRedisKey(item.App, item.AccessKey)
	return nil
}

func (s *apiKeyService) _findByAk(app string, accessKey string) (*xapikey.Aksk, error) {
	filter := bson.M{
		"app":       app,
		"accessKey": accessKey,
	}
	return s.FindOne(filter)
}

func (s *apiKeyService) _deleteRedisKey(app string, ak string) {
	key := getRedisKeyForXApiKey(app, ak)
	err := getServiceGroup().redisService.DeleteKey(key)
	if err != nil {
		log.Logger.Warn(fmt.Sprintf("在删除xapikey的redis key时出现异常,key:%s,err:%s",
			key,
			err.Error()))
	}
}

func getRedisKeyForXApiKey(app string, ak string) string {
	return fmt.Sprintf("ak:%s:%s", app, ak)
}
