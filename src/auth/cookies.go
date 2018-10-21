package auth

import (
	"time"
	"net/http"
)

type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	Raw      string
	Unparsed []string
}

func SetCookie(rw http.ResponseWriter, tokenName string, token string) string {

	expiration := time.Now().Add(365 * 24 * time.Hour)

	cookie := http.Cookie{Name: tokenName, Value: token, Expires: expiration}

	http.SetCookie(rw, &cookie)

	return token

}
