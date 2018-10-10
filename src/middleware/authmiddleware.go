package middleware

import (
	"net/http"
	"fmt"
)

const (
	ErrTokenInvalid = "Token is invalid"
	ErrTokenMissing = "No token is present"
	ErrTokenExpired = "Token is expired"
	ErrAccessDenied = "Access to this resource is denied"
)


var cookie http.Cookie

func AuthMiddleware(rw http.ResponseWriter, r *http.Request, nextFunc http.HandlerFunc) {

/* 1) Check for existing token. If not found, fetch new one */

	existingToken, _ := r.Cookie("auth_token")

	fmt.Println(existingToken)

	/*jwt, err := jwt.Parse(existingToken, func(token *jwt.Token) (interface{}, error) {
		fmt.Println(error)
	})*/


	/* 3) Check access against application */

	/* 4) Execute request */
	//nextFunc(rw, r)

	/* 5) Issue refresh token */
	//fmt.Println("auth middleware -> after the controller was executed")

}
