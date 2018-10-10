package auth

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"../config/mainconf"
	"net/url"
	"log"
	"fmt"
	"time"
)

var AuthConf = mainconf.GetAuthConfig()
var Conf = mainconf.BuildConfig()

func AuthCallback(w http.ResponseWriter, r *http.Request) {

	authCode := r.URL.Query().Get("code")

	AccToken, err := FetchAccessToken(authCode)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(AccToken)

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


func Login(w http.ResponseWriter, r *http.Request){

	url := AuthConf.AuthHost+"/authorize?client_id="+AuthConf.AuthClientId+"&response_type=code&response_mode=query&scope=openid&state=12345&redirect_uri="+Conf.ApiHost+"/authcallback"

	http.Redirect(w, r, url, http.StatusSeeOther)

}


func Logout(w http.ResponseWriter, r *http.Request){

	c := &http.Cookie{
		Name: "auth_token",
		Value:    "",
		Path:     "/",
		Expires: time.Unix(0, 0),
		HttpOnly: true,
	}

	http.SetCookie(w, c)

	url := AuthConf.AuthHost+"/logout?post_logout_redirect_uri="+Conf.ApiHost

	http.Redirect(w, r, url, http.StatusSeeOther)

}

func ValidateToken(){

}

func RefreshToken(){

}
