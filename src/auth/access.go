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
)


var conf = mainconf.BuildConfig()

// Instantiate access control handlers
var gormAdapter = gormadapter.NewAdapter("mssql", "sqlserver://"+conf.SqlUser+":"+url.QueryEscape(conf.SqlPass)+"@"+conf.SqlHost+":"+strconv.Itoa(conf.SqlPort)+"?database="+conf.SqlDB, true)
var AccessEnforcer = casbin.NewEnforcer("access_control_model.conf", gormAdapter)

/**
 * Pulls access keys from Azure endpoint and caches them in database
 */
func CacheAccessKeys(cacheMethod string){

	keyMap, err := RetrieveAccessKeys()
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

		database.ExecuteRawQuery(sqlUpdate)


	} else if cacheMethod == "local_env" {

		jsonKeys := formatKeysAsJson(keyMap)
		os.Setenv("ACCESS_KEYS", jsonKeys)

	} else if cacheMethod == "local_file" {

		//jsonKeys := formatKeysAsJson(keyMap)
		//os.Setenv("ACCESS_KEYS", jsonKeys)

	} else if cacheMethod == "memory" {

	}
}

func formatKeysAsJson(keyMap map[string][]map[string]interface{}) string {
	jsonKeys := "["
	for _, val := range keyMap {
		for _, v := range val {
			jsonKeys += "{'kid': " + v["kid"].(string) + ", 'x5c': " + v["x5c"].([]interface{})[0].(string) + "},"
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


func RetrieveAccessKeys() (map[string][]map[string]interface{}, error) {
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
				CacheAccessKeys(cacheMethod)
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


func GetAccessKey(kid string) (string, error) {

	var kidMap = make(map[string]string)
	kidMap["kid"] = kid

	sql := database.Statement{ Sql: `SELECT x5c FROM dbo.access_keys WHERE kid = {{kid}}`, Params: kidMap}

	sqlResult, err := database.SelectAndReturnResultSet(sql)

	result := sqlResult[0]["x5c"]

	return result.(string), err

}

func ValidateToken(cookie *http.Cookie) bool {

	var authKey = `-----BEGIN CERTIFICATE-----\n`
	claims := jwt.MapClaims{}


	tkn, _ := jwt.Parse(cookie.Value, func(tkn *jwt.Token) (interface{}, error) {
		return tkn, nil
	})

	jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error){
		return []byte(AuthConf.AuthSecret), nil
	})

	// Gets key ID from JWT for validating against openid server
	kid := tkn.Header["kid"]

	cachedKey, err := GetAccessKey(kid.(string))

	if err != nil {
		util.ErrorHandler(err)
	}

	authKey += cachedKey

	authKey += `\n-----END CERTIFICATE-----\n`

	tkn.SignedString(authKey)

	valid := tkn.Valid

	if valid != true {



	}

	return valid
}

func ProtectedEndpoint (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	cookie, err := r.Cookie("id_token")
	if err != nil {
		util.ErrorHandler(err)
		w.Write([]byte(err.Error()))
		return
	}
	validToken := ValidateToken(cookie)

	fmt.Println(validToken)

	return

	// Return from Token Decoding method here
	/*claims := jwt.MapClaims{}

	cookie, err := r.Cookie("id_token")
	if err != nil {
		util.ErrorHandler(err)
		w.Write([]byte(err.Error()))
		return
	}

	//w.Write([]byte(cookie.Value))
	//return

	jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error){
		return []byte(AuthConf.AuthSecret), nil
	})


	username := claims["unique_name"]

	accRes := AccessEnforcer.Enforce(username, r.URL.String(), r.Method)
	if accRes != false {
		next(w, r)
	} else {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("User not authorized to access this resource")
		return
	}*/

	/*cookie, _ := r.Cookie("id_token")

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}


		//signed, err := token.SignedString(key)


		return signed, err
	})



	fmt.Println(token.Valid)

	if err != nil {
		util.ErrorHandler(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user User
		mapstructure.Decode(claims, &user)
		json.NewEncoder(w).Encode(user)
	} else {
		json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
	}



	/*if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user User
		mapstructure.Decode(claims, &user)

		json.NewEncoder(w).Encode(user)







	} else {
		w.Write([]byte(err.Error()))
		json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
	}*/

}




/********** DEPRECATED ***********
 * Checks if user has access to given method on a database object
 * `objectName` is fully qualified, e.g. Schema.Table
 */
func CheckUserAccess(userName string, objectName string, method string) (int, error){

	params := map[string]string{"userName": userName, "objectName": objectName, "method": method}

	sql := database.Statement{ Sql: `SELECT count(ur.id) as result_count
  										FROM People.Users u
  										INNER JOIN Security.User_Roles ur on ur.user_id = u.id
  										INNER JOIN Security.Roles r on r.id = ur.role_id
										INNER JOIN Security.Object_Role_Permissions orp on orp.role_id = r.id
  										INNER JOIN Security.Permissions p on p.id = orp.permission_id
  										INNER JOIN Security.Objects o on o.id = orp.object_id
									WHERE u.username = '{{userName}}'
										AND o.schema_name + '.' + o.object_name = '{{objectName}}'
  										AND p.name = '{{method}}'`, Params: params }

	result, err := database.SelectAndReturnCount(sql)
	if err != nil {
		util.ErrorHandler(err)
	}

	return result, err
}