package auth

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"../config/mainconf"
	"net/url"
	"log"
	"time"
	"fmt"
)

var AuthConf = mainconf.GetAuthConfig()
var Conf = mainconf.BuildConfig()

func AuthCallback(w http.ResponseWriter, r *http.Request) {

	/*authCode := r.URL.Query().Get("code")

	AccToken, err := FetchAccessToken(authCode)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(AccToken)*/

	authErr := r.URL.Query().Get("error")

	fmt.Println(authErr)

	if authErr != "" {
		authErrDesc := r.URL.Query().Get("error_description")
		fmt.Println(authErrDesc)
	}

	IdToken := r.URL.Query().Get("id_token")

	// Replace this with access token logic
	AccToken := IdToken

	SetCookie(w, "auth_token", AccToken)

}

func FetchAccessToken(auth_code string) (string, error){

	tokenUrl := AuthConf.AuthHost+"/token"

	resp, err := http.PostForm(tokenUrl, url.Values{
		"client_id": {AuthConf.AuthClientId},
		"client_secret": {AuthConf.AuthSecret},
		"code": {auth_code},
		"scope": {"https://graph.microsoft.com/.default"},
		"grant_type": {"client_credentials"},
	})

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	rsBody, err := ioutil.ReadAll(resp.Body)

	type WithAccToken struct {
		AccessToken string `json:"access_token"`
	}

	var dat WithAccToken

	err = json.Unmarshal(rsBody, &dat)
	if err != nil {
		return "",err
	}

	return dat.AccessToken, err
}


func LoginB2C(w http.ResponseWriter, r *http.Request){

	redirUrl := r.URL.Query().Get("redirect_uri")

	url := AuthConf.AuthHost+"/authorize?p=B2C_1_SignUpIn&client_id="+AuthConf.AuthClientId+"&nonce=defaultNonce&redirect_uri="+redirUrl+"&scope=https%3A%2F%2FFortisureB2CTenant.onmicrosoft.com%2Fapi-dev%2Fadmin&response_type=token&prompt=login"

	http.Redirect(w, r, url, http.StatusSeeOther)

}

func LoginOrg(w http.ResponseWriter, r *http.Request){

	redirUrl := r.URL.Query().Get("redirect_uri")

	url := AuthConf.AuthHost+"/authorize?client_id="+AuthConf.AuthClientId+"&redirect_uri="+redirUrl+"&scope=https%3A%2F%2Fgraph.microsoft.com%2Fme&response_type=id_token&nonce=defaultNonce"

	http.Redirect(w, r, url, http.StatusSeeOther)

}

func LogoutOrg(w http.ResponseWriter, r *http.Request){

	c := &http.Cookie{
		Name: "auth_token",
		Value:    "",
		Path:     "/",
		Expires: time.Unix(0, 0),
		HttpOnly: true,
	}

	http.SetCookie(w, c)

	url := AuthConf.AuthHost+"/logout?client_id="+AuthConf.AuthClientId+"&post_logout_redirect_uri="+Conf.ApiHost

	http.Redirect(w, r, url, http.StatusSeeOther)

}

func LogoutB2C(w http.ResponseWriter, r *http.Request){

	c := &http.Cookie{
		Name: "auth_token",
		Value:    "",
		Path:     "/",
		Expires: time.Unix(0, 0),
		HttpOnly: true,
	}

	http.SetCookie(w, c)

	url := AuthConf.AuthHost+"/logout?p=B2C_1_SignUpIn&post_logout_redirect_uri="+Conf.ApiHost

	http.Redirect(w, r, url, http.StatusSeeOther)

}

func ValidateToken(w http.ResponseWriter, r *http.Request){

	fmt.Println(r.URL.Query().Get("id_token"))

}

func RefreshToken(){

}
