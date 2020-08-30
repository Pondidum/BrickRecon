package preen

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
)

type UserInfo struct {
	Name          string
	Authenticated bool
}

type AuthMiddleware struct {
	Username string
	Password string
	Realm    string
}

func NewBasicAuthMiddlware(user, pass, realm string) *AuthMiddleware {
	return &AuthMiddleware{
		Username: user,
		Password: pass,
		Realm:    realm,
	}
}

func (a *AuthMiddleware) authenticate(user, pass string, r *http.Request) bool {
	// Equalize lengths of supplied and required credentials by hashing them
	givenUser := sha256.Sum256([]byte(user))
	givenPass := sha256.Sum256([]byte(pass))
	requiredUser := sha256.Sum256([]byte(a.Username))
	requiredPass := sha256.Sum256([]byte(a.Password))

	// Compare the supplied credentials to those set in our options
	if subtle.ConstantTimeCompare(givenUser[:], requiredUser[:]) == 1 &&
		subtle.ConstantTimeCompare(givenPass[:], requiredPass[:]) == 1 {
		return true
	}

	return false
}

func (a *AuthMiddleware) Middleware(c *MiddlewareContext, request *http.Request, response http.ResponseWriter) bool {

	givenUser, givenPass, given := request.BasicAuth()

	isAuthenticated := given && a.authenticate(givenUser, givenPass, request)

	user := UserInfo{
		Name:          givenUser,
		Authenticated: isAuthenticated,
	}

	context.Set(request, "UserInfo", user)

	if c.AuthRequired() || request.Method == http.MethodPost {

		if user.Authenticated == false {

			response.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm=%q`, a.Realm))
			http.Error(response, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return false
		}

	}

	return true
}
