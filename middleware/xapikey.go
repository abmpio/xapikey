package middleware

import (
	"fmt"
	"strings"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/app/host"
	"github.com/abmpio/irisx/controllerx"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"

	"github.com/kataras/iris/v12/core/router"
)

func UseXApiKey(apiBuilder *router.APIBuilder) {
	apiBuilder.Use(serveHTTP)
}

func serveHTTP(ctx *context.Context) {
	claims := controllerx.GetCasdoorMiddleware().GetUserClaims(ctx)
	if claims != nil {
		ctx.Next()
		return
	}
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
	appName := host.GetHostEnvironment().GetEnvString(host.ENV_AppName)
	if len(appName) <= 0 {
		ctx.Next()
		return
	}
	apiKey, err := getServiceGroup().apikeyService.FindByAk(appName, ak)
	if err != nil {
		log.Logger.Warn(fmt.Sprintf("find x-api-key error: %v", err))
		ctx.StopExecution()
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString(err.Error())
		return
	}
	if apiKey == nil {
		err := fmt.Errorf(fmt.Sprintf("invalid x-api-key,ak:%s", ak))
		log.Logger.Warn(err.Error())
		ctx.StopExecution()
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString(err.Error())
		return
	}
	if apiKey.SecretKey != sk {
		err := fmt.Errorf(fmt.Sprintf("invalid x-api-key,ak:%s", ak))
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
