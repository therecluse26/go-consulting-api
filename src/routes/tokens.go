package routes

import (
	"net/http"
	"time"
	"github.com/dgrijalva/jwt-go"
	"os"
)

func SetTokenData(token jwt.Token) {

	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["admin"] = true
	claims["name"] = "Brad Magyar"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

}

func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	//authConf := mainconf.GetAuthConfig()

	SetTokenData(*token)

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString([]byte(os.Getenv("AuthSecret")))

	/* Finally, write the token to the browser window */
	w.Write([]byte(tokenString))
}
