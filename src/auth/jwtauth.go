package auth

import (
	"net/http"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"bytes"
	"io/ioutil"
	"../util"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

func DecodeJWT(jwt_token string)  {

	jwt.DecodeSegment("payload")

}

func GetTokenScope(tokUrl string, clientId string, secret string) (string,error){

	body := bytes.NewBuffer([]byte("grant_type=client_credentials&client_id="+clientId+"&client_secret="+secret+"&response_type=token"))

	req, err := http.NewRequest("POST",tokUrl,body)

	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		util.ErrorHandler(err)
	}

	defer resp.Body.Close()

	rsBody, err := ioutil.ReadAll(resp.Body)

	type WithScope struct {
		Scope string `json:"scope"`
	}

	var dat WithScope

	err = json.Unmarshal(rsBody,&dat)
	if err != nil {
		util.ErrorHandler(err)
	}

	return dat.Scope,err
}