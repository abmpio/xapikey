package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/app/host"
	"github.com/abmpio/entity/tenancy"
	"github.com/abmpio/irisx/controllerx"
	webapp "github.com/abmpio/webserver/app"
	"github.com/abmpio/webserver/controller"
	"github.com/abmpio/xapikey"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	c.EntityController.Options.EnableFilterCurrentUser = true
	log.Logger.Debug(fmt.Sprintf("正在构建路由,%s...", routerPath))

	routerParty := c.EntityController.RegistRouter(webapp,
		controllerx.BaseEntityControllerWithRouterPath(routerPath),
		controllerx.BaseEntityControllerWithAllEndpointDisabled(true),
	)
	routerParty.Post("/", c.create)
	routerParty.Get("/all", c.all)
	routerParty.Delete("/{id}", c.delete)
}

type xapiKeyCreateInput struct {
	// 所属app
	Alias          string     `json:"alias"`
	Description    string     `json:"description"`
	ExpirationTime *time.Time `json:"expirationTime"`
	Status         bool       `json:"status"`
	IpWhitelist    string     `json:"ipWhitelist"`
}

func (c *apiKeyController) create(ctx iris.Context) {
	input := &xapiKeyCreateInput{}
	err := ctx.ReadJSON(&input)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, err)
		return
	}
	if len(input.Alias) <= 0 {
		controller.HandleErrorBadRequest(ctx, fmt.Errorf("alias 字段不能为空"))
		return
	}
	app := host.GetHostEnvironment().GetEnvString(host.ENV_AppName)
	if len(app) <= 0 {
		controller.HandleErrorInternalServerError(ctx, fmt.Errorf("本应用不支持创建api key"))
		return
	}
	tenantId := tenancy.TenantIdFromContext(ctx)
	newAksk := &xapikey.Aksk{
		App:            app,
		Alias:          input.Alias,
		Status:         input.Status,
		IpWhitelist:    input.IpWhitelist,
		Description:    input.Description,
		ExpirationTime: input.ExpirationTime,
	}
	// set current tenant
	newAksk.TenantId = tenantId
	if newAksk.ExpirationTime != nil {
		newAksk.ExpirationTime = input.ExpirationTime
	}
	ak, sk := xapikey.GenerateAKSK()
	newAksk.AccessKey = ak
	newAksk.SecretKey = sk

	// handler user info
	c.SetUserInfo(ctx, newAksk)

	newItem, err := getServiceGroup().apikeyService.Create(newAksk)
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

// delete
func (c *apiKeyController) delete(ctx iris.Context) {
	idValue := ctx.Params().Get("id")
	if len(idValue) <= 0 {
		controller.HandleErrorBadRequest(ctx, errors.New("id must not be empty"))
		return
	}
	oid, err := primitive.ObjectIDFromHex(idValue)
	if err != nil {
		controller.HandleErrorBadRequest(ctx, fmt.Errorf("invalid id format,err:%s", err.Error()))
		return
	}
	err = getServiceGroup().apikeyService.Delete(oid)
	if err != nil {
		controller.HandleErrorInternalServerError(ctx, err)
		return
	}
	controller.HandleSuccess(ctx)
}
