package preen

import (
	"fmt"
	"net/http"
	"strings"
)

func TurbolinksMiddleware(c *MiddlewareContext, request *http.Request, response http.ResponseWriter) bool {
	// ctx := request.Context()

	reqType := request.Header.Get("X-Requested-With")

	if strings.EqualFold(reqType, "XMLHttpRequest") == false {
		return true
	}

	script := fmt.Sprintf(`
Turbolinks.clearCache()
Turbolinks.visit("%s", { action: "replace" })
`, request.RequestURI)

	response.Header().Add("Content-Type", "text/javascript")
	response.Header().Add("X-Xhr-Redirect", request.RequestURI)
	response.Write([]byte(script))

	return false
}
