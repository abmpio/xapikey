package middleware

import (
	"fmt"
	"strings"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/app/host"
	"github.com/abmpio/entity/tenancy"
	"github.com/abmpio/irisx/controllerx"
	"github.com/abmpio/xapikey"
	"github.com/abmpio/xapikey/options"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"

	"github.com/kataras/iris/v12/core/router"
)

func UseXApiKey(apiBuilder *router.APIBuilder) {
	apiBuilder.Use(serveHTTP)
}

func serveHTTP(ctx *context.Context) {
	// claims := controllerx.GetCasdoorMiddleware().GetUserClaims(ctx)
	// if claims != nil {
	// 	ctx.Next()
	// 	return
	// }
	ak, sk, err := extractXApiKey(ctx)
	if err != nil {
		log.Logger.Warn(fmt.Sprintf("Error extracting x-api-key: %v", err))
		ctx.StopExecution()
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString(err.Error())
		return
	}
	if len(ak) <= 0 || len(sk) <= 0 {
		// no x-api-key
		ctx.Next()
		return
	}
	appNameList := make([]string, 0)
	appName := host.GetHostEnvironment().GetEnvString(host.ENV_AppName)
	if len(appName) > 0 {
		appNameList = append(appNameList, appName)
	}
	extraAppName := options.GetOptions().ExtraAppName
	if len(extraAppName) > 0 {
		appNameList = append(appNameList, extraAppName)
	}
	if len(appNameList) <= 0 {
		ctx.Next()
		return
	}
	var apiKey *xapikey.Aksk
	for _, eachAppName := range appNameList {
		apiKey, err = getServiceGroup().apikeyService.FindByAk(tenancy.TenantIdFromContext(ctx), eachAppName, ak)
		if err != nil {
			log.Logger.Warn(fmt.Sprintf("find x-api-key error: %v", err))
			ctx.StopExecution()
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.WriteString(err.Error())
			return
		}
		if apiKey != nil {
			break
		}
	}
	if apiKey == nil {
		err := fmt.Errorf("invalid x-api-key,ak:%s", ak)
		log.Logger.Warn(err.Error())
		ctx.StopExecution()
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString(err.Error())
		return
	}
	if apiKey.SecretKey != sk {
		err := fmt.Errorf("invalid x-api-key,ak:%s", ak)
		log.Logger.Warn(err.Error())
		ctx.StopExecution()
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString(err.Error())
		return
	}
	expired := apiKey.CheckExpired()
	if expired {
		// 已经过期
		err := fmt.Errorf("invalid x-api-key,ak:%s", ak)
		log.Logger.Warn(err.Error())
		ctx.StopExecution()
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString(err.Error())
		return
	}
	claim := &casdoorsdk.Claims{
		User: casdoorsdk.User{
			Id: apiKey.CreatorId,
		},
	}
	// set claims
	ctx.Values().Set(controllerx.GetCasdoorMiddleware().Options.Jwt.ContextKey, claim)
	ctx.Values().Set("userId", claim.Id)

	ctx.Next()
}

func extractXApiKey(ctx iris.Context) (string, string, error) {
	xapiKey := ctx.GetHeader("x-api-key")
	if xapiKey == "" {
		return "", "", nil
	}
	headerParts := strings.Split(xapiKey, " ")
	if len(headerParts) != 2 {
		return "", "", fmt.Errorf("x-api-key format must be {x-api-ak} {x-api-sk}")
	}
	return headerParts[0], headerParts[1], nil
}
