package auth

import (
	"net/http"
	"github.com/casbin/gorm-adapter"
	"github.com/casbin/casbin"
	"../config/mainconf"
	"../database"
	"strconv"
	"net/url"
	"time"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"encoding/json"
	"../util"
	"os"
	"io/ioutil"
)


var conf = mainconf.BuildConfig()

// Instantiate access control handlers
var gormAdapter = gormadapter.NewAdapter("mssql", "sqlserver://"+conf.SqlUser+":"+url.QueryEscape(conf.SqlPass)+"@"+conf.SqlHost+":"+strconv.Itoa(conf.SqlPort)+"?database="+conf.SqlDB, true)
var AccessEnforcer = casbin.NewEnforcer("access_control_model.conf", gormAdapter)


func ProtectedEndpoint (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	cookie, err := r.Cookie("id_token")

	fmt.Println(cookie)

	if err != nil {
		util.ErrorHandler(err)
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}

	validToken, err := ValidateToken(cookie)
	if !validToken {
		util.ErrorHandler(err)
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}

	accRes := verifyUserPermissions(cookie, r.URL.String(), r.Method)
	if accRes != false {
		next(w, r)
	} else {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("User not authorized to access this resource")
		return
	}

}


/**
 * Pulls access keys from Azure endpoint and caches them in database
 */
func cacheAccessKeys(cacheMethod string){

	keyMap, err := retrieveAccessKeysFromAzure()
	if err != nil {
		util.ErrorHandler(err)
		return
	}

	if cacheMethod == "database" {

		sqlUpdate := `BEGIN TRY
						BEGIN TRANSACTION
							DELETE FROM dbo.access_keys;
    						INSERT INTO dbo.access_keys (kid, x5c) VALUES`

		for _, value := range keyMap["keys"] {
			sqlUpdate += `('`+value["kid"].(string)+`', '`+value["x5c"].([]interface{})[0].(string)+`'),`
		}

		ln := len(sqlUpdate)
		// Trims last comma after values string
		if sqlUpdate[ln-1] == ',' {
			sqlUpdate = sqlUpdate[:ln-1]
			sqlUpdate += `;`
		}
		sqlUpdate += `COMMIT TRANSACTION
					END TRY
					BEGIN CATCH
						IF @@TRANCOUNT > 0
						ROLLBACK TRAN
					END CATCH`

		err := database.ExecuteRawQuery(sqlUpdate)
		if err != nil {
			util.ErrorHandler(err)
		}


	} else if cacheMethod == "local_env" {

		jsonKeys := formatKeysAsJson(keyMap)
		os.Setenv("ACCESS_KEYS", jsonKeys)

	} else if cacheMethod == "local_file" {

		jsonKeys := formatKeysAsJson(keyMap)
		ioutil.WriteFile("access_keys", []byte(jsonKeys), 0664)

	} else if cacheMethod == "memory" {

		jsonKeys := formatKeysAsJson(keyMap)
		database.MemcachedStore("access_keys", jsonKeys)

	}
}

func RetrieveLocalAccessKeys(cacheMethod string) (string, error) {

	var result string
	var err error

	if cacheMethod == "database" {

		sql := database.Statement{Sql: "SELECT kid, x5c from dbo.access_keys"}
		sqlResult, e := database.SelectAndReturnResultSet(sql)
		if e != nil {
			err = e
		}

		res, e := json.Marshal(sqlResult)
		if err != nil {
			err = e
		}

		result = string(res)

	} else if cacheMethod == "local_env" {

		result = os.Getenv("ACCESS_KEYS")
		err = nil

	} else if cacheMethod == "local_file" {

		res, fileErr := ioutil.ReadFile("access_keys")

		result = string(res)
		err = fileErr


	} else if cacheMethod == "memory" {

		result, err = database.MemcachedRetrieve("access_keys")
	}

	return result, err

}

func formatKeysAsJson(keyMap map[string][]map[string]interface{}) string {
	jsonKeys := "["
	for _, val := range keyMap {
		for _, v := range val {
			jsonKeys += `{"kid": ` + `"` + v["kid"].(string) + `"` + `, "x5c": ` + `"` + v["x5c"].([]interface{})[0].(string) + `"` + `},`
		}
	}
	ln := len(jsonKeys)
	// Trims last comma after values string
	if jsonKeys[ln-1] == ',' {
		jsonKeys = jsonKeys[:ln-1]
	}
	jsonKeys += "]"

	return jsonKeys
}


func retrieveAccessKeysFromAzure() (map[string][]map[string]interface{}, error) {
	var openIdConf map[string]string
	var keyMap map[string][]map[string]interface{}

	// Pulls keys from openid server
	openIdConfRaw, err := http.Get("https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration")
	if err != nil {
		util.ErrorHandler(err)
	}
	json.NewDecoder(openIdConfRaw.Body).Decode(&openIdConf)
	jwksUri := openIdConf["jwks_uri"]
	keyMapRaw, _ := http.Get(jwksUri)

	json.NewDecoder(keyMapRaw.Body).Decode(&keyMap)

	return keyMap, err
}

/**
 * Refreshes access policies from database every x seconds
 */
func CacheAccessKeysTimer(seconds time.Duration, cacheMethod string) {

	timer := time.NewTicker(seconds * time.Second)

	go func() {
		for {
			select {
			case <- timer.C:
				cacheAccessKeys(cacheMethod)
				fmt.Println("Refreshed Access Keys")
			}
		}
	}()
}


/**
 * Refreshes access policies from database every x seconds
 */
func LoadAccessPolicyLoopTimer(seconds time.Duration) {

	timer := time.NewTicker(seconds * time.Second)

	go func() {
		for {
			select {
				case <- timer.C:
				AccessEnforcer.LoadPolicy()
			}
		}
	}()
}


func getAccessKey(kid string) (string) {

	var decodedKeys []map[string]string
	var authKey = "-----BEGIN CERTIFICATE-----\n"

	keys, err := RetrieveLocalAccessKeys(conf.CacheMethod)
	if err != nil {

		util.ErrorHandler(err)

		// Fetch remote key
		remoteKeyMap, azerr := retrieveAccessKeysFromAzure()
		if azerr != nil {
			util.ErrorHandler(azerr)
		}
		for _, value := range remoteKeyMap["keys"] {
			if value["kid"].(string) == kid {
				authKey += value["x5c"].([]interface{})[0].(string)
			}
		}


	} else {

		// Fetch local key
		json.Unmarshal([]byte(keys), &decodedKeys)
		for _, val := range decodedKeys {
			if val["kid"] == kid {
				authKey += val["x5c"]
			}
		}

	}

	authKey += "\n-----END CERTIFICATE-----\n"

	return authKey

}

func ValidateToken(cookie *http.Cookie) (bool, error) {

	tkn, err := jwt.Parse(cookie.Value, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Error parsing id token")
		}
		valid, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(getAccessKey(tkn.Header["kid"].(string))))
		return valid, nil
	})

	return tkn.Valid, err
}

func verifyUserPermissions (cookie *http.Cookie, url string, method string) bool {
	claims := jwt.MapClaims{}

	jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error){
		return []byte(AuthConf.AuthSecret), nil
	})

	username := claims["unique_name"]

	accRes := AccessEnforcer.Enforce(username, url, method)

	return accRes
}
