package preen

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
)

type basicAuth struct {
	realm    string
	user     string
	password string
}

type AuthOptions struct {
	Realm    string
	User     string
	Password string
}

func BasicAuthMiddleware(o AuthOptions) basicAuth {
	return basicAuth{
		user:     o.User,
		password: o.Password,
		realm:    o.Realm,
	}
}

func (b *basicAuth) UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		givenUser, givenPass, given := r.BasicAuth()
		isAuthenticated := given && b.simpleBasicAuth(givenUser, givenPass, r)

		context.Set(r, "UserInfo", UserInfo{
			Name:          givenUser,
			Authenticated: isAuthenticated,
		})

		next.ServeHTTP(w, r)
	})
}

func (b *basicAuth) Wrap(next ViewMiddleware) ViewMiddleware {
	return func(w http.ResponseWriter, r *http.Request, model interface{}) {

		user := context.Get(r, "UserInfo").(UserInfo)

		if user.Authenticated == false {
			b.requestAuth(w, r)
			return
		}

		next(w, r, model)
	}
}

func (b *basicAuth) simpleBasicAuth(user, pass string, r *http.Request) bool {
	// Equalize lengths of supplied and required credentials by hashing them
	givenUser := sha256.Sum256([]byte(user))
	givenPass := sha256.Sum256([]byte(pass))
	requiredUser := sha256.Sum256([]byte(b.user))
	requiredPass := sha256.Sum256([]byte(b.password))

	// Compare the supplied credentials to those set in our options
	if subtle.ConstantTimeCompare(givenUser[:], requiredUser[:]) == 1 &&
		subtle.ConstantTimeCompare(givenPass[:], requiredPass[:]) == 1 {
		return true
	}

	return false
}

func (b *basicAuth) requestAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm=%q`, b.realm))

	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
