package xapikey

import (
	"github.com/kataras/iris/v12/context"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

type IXApiUserInfoService interface {
	// set user info
	SetupUserInfo(ctx *context.Context, u *casdoorsdk.Claims) error
}
