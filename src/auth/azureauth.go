package auth

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/dgrijalva/jwt-go"
	"bytes"
	"io/ioutil"
)

func Authenticate(){


}

func GetToken(){

}

func asdf(authHost string, clientId string, secret string, scope string, refresh_token string) (string,error){

	body := bytes.NewBuffer([]byte("grant_type=refresh_token&client_id="+clientId+"&client_secret="+secret+"&scope="+scope+"&refresh_token="+refresh_token))

	resp, err := http.Post(authHost+"/token", "application/x-www-form-urlencoded", body)

	fmt.Println(resp.Status)

	if err != nil {
		return "",err
	}

	defer resp.Body.Close()

	rsBody, err := ioutil.ReadAll(resp.Body)

	type WithAccToken struct {
		AccessToken string `json:"access_token"`
	}

	var dat WithAccToken

	err = json.Unmarshal(rsBody,&dat)
	if err != nil {
		return "",err
	}

	return dat.AccessToken,err
}


func ProtectedEndpoint (w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	token, _ := jwt.Parse(params["token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user User
		mapstructure.Decode(claims, &user)

		json.NewEncoder(w).Encode(user)
	} else {
		json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
	}

}

func GetTokenScope(tokUrl string, clientId string, secret string) (string,error){

	body := bytes.NewBuffer([]byte("grant_type=client_credentials&client_id="+clientId+"&client_secret="+secret+"&response_type=token"))

	req, err := http.NewRequest("POST",tokUrl,body)

	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "",err
	}

	defer resp.Body.Close()

	rsBody, err := ioutil.ReadAll(resp.Body)

	type WithScope struct {
		Scope string `json:"scope"`
	}

	var dat WithScope

	err = json.Unmarshal(rsBody,&dat)
	if err != nil {
		return "",err
	}

	return dat.Scope,err
}