package controller

import (
	"fmt"
	"strings"

	"github.com/abmpio/xapikey/options"

	webapp "github.com/abmpio/webserver/app"
)

func registControllers(webApp *webapp.Application) {
	newApiKeyController().RegistRouter(webApp, appendRouterPrefixPath("/xapikey"))
}

func appendRouterPrefixPath(path string) string {
	o := options.GetOptions()
	prefixPath := o.RouterPrefixPath
	prefixPath = strings.TrimPrefix(prefixPath, "/")
	prefixPath = strings.TrimSuffix(prefixPath, "/")
	if len(prefixPath) <= 0 {
		return path
	}
	path = strings.TrimPrefix(path, "/")
	if len(path) <= 0 {
		return fmt.Sprintf("/%s", o.RouterPrefixPath)
	} else {
		return fmt.Sprintf("/%s/%s", o.RouterPrefixPath, path)
	}
}
