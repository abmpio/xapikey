package controller

import (
	"fmt"
	"time"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/irisx/controllerx"
	"github.com/abmpio/mongodbr"
	webapp "github.com/abmpio/webserver/app"
	"github.com/abmpio/webserver/controller"
	"github.com/abmpio/xapikey"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
)

type apiKeyController struct {
	controllerx.EntityController[xapikey.Aksk]
}

func newApiKeyController() *apiKeyController {
	c := &apiKeyController{}
	c.EntityController = controllerx.EntityController[xapikey.Aksk]{
		EntityService: getServiceGroup().apikeyService,
	}
	return c
}

func (c *apiKeyController) RegistRouter(webapp *webapp.Application, routerPath string) {
	c.EntityController.RouterPath = routerPath
	c.EntityController.Options.EnableFilterCurrentUser = true
	log.Logger.Debug(fmt.Sprintf("正在构建路由,%s...", routerPath))

	routerParty := c.EntityController.RegistRouter(webapp,
		controllerx.BaseEntityControllerWithAllEndpointDisabled(true),
	)
	routerParty.Post("/", c.MergeAuthenticatedContextIfNeed(c.Options.AuthenticatedDisabled, c.create)...)
	routerParty.Get("/all", c.MergeAuthenticatedContextIfNeed(c.Options.AuthenticatedDisabled, c.all)...)
}

type xapiKeyCreateInput struct {
	// 所属app
	App            string     `json:"app"`
	Alias          string     `json:"alias"`
	Description    string     `json:"description"`
	ExpirationTime *time.Time `json:"expirationTime"`
	Status         bool       `json:"status"`
	IpWhitelist    string     `json:"ipWhitelist"`

	Properties map[string]interface{} `json:"properties"`
}

func (c *apiKeyController) create(ctx iris.Context) {
	currentUserId := controllerx.GetUserId(ctx)
	input := &xapiKeyCreateInput{}
	err := ctx.ReadJSON(&input)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, err)
		return
	}
	newXapiKey := &xapikey.Aksk{
		App:            input.App,
		Alias:          input.Alias,
		Status:         input.Status,
		IpWhitelist:    input.IpWhitelist,
		Description:    input.Description,
		ExpirationTime: input.ExpirationTime,
	}
	if newXapiKey.ExpirationTime != nil {
		newXapiKey.ExpirationTime = input.ExpirationTime
	}
	ak, sk := xapikey.GenerateAKSK()
	newXapiKey.AccessKey = ak
	newXapiKey.SecretKey = sk
	err = mongodbr.Validate(newXapiKey)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, err)
		return
	}

	newXapiKey.BeforeCreate()
	// handler user info
	c.SetUserInfo(ctx, newXapiKey)

	list, err := getServiceGroup().apikeyService.FindList(bson.M{
		"creatorId": currentUserId,
		"alias":     newXapiKey.Alias,
	})
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	if len(list) > 0 {
		controller.HandleErrorInternalServerError(ctx, fmt.Errorf("已经存在着名称为: %s 的其它记录", input.Alias))
		return
	}
	newItem, err := c.GetEntityService().Create(newXapiKey)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccessWithData(ctx, newItem)
}

func (c *apiKeyController) all(ctx iris.Context) {
	var input struct {
		App string `url:"app"`
	}
	err := ctx.ReadQuery(&input)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, err)
		return
	}
	if len(input.App) <= 0 {
		controller.HandleErrorBadRequest(ctx, fmt.Errorf("app参数不能为空"))
		return
	}
	filter := map[string]interface{}{
		"app": input.App,
	}
	if c.Options.EnableFilterCurrentUser {
		// auto filter current userId
		controllerx.AddUserIdFilterIfNeed(filter, &xapikey.Aksk{}, ctx)
	}

	if c.Options.ListFilterFunc != nil {
		c.Options.ListFilterFunc(&xapikey.Aksk{}, filter, ctx)
	}
	var list []*xapikey.Aksk
	if len(filter) > 0 {
		list, err = c.GetEntityService().FindList(filter)
	} else {
		list, err = c.GetEntityService().FindAll()
	}
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccessWithData(ctx, list)
}
