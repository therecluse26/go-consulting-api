package auth

import (
	"net/http"
	"github.com/casbin/gorm-adapter"
	"github.com/casbin/casbin"
	"../config/mainconf"
	"../database"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"net/url"
	"time"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

var conf = mainconf.BuildConfig()

// Instantiate access control handlers
var gormAdapter = gormadapter.NewAdapter("mssql", "sqlserver://"+conf.SqlUser+":"+url.QueryEscape(conf.SqlPass)+"@"+conf.SqlHost+":"+strconv.Itoa(conf.SqlPort)+"?database="+conf.SqlDB, true)
var AccessEnforcer = casbin.NewEnforcer("access_control_model.conf", gormAdapter)



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


func ProtectedEndpoint (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// Return from Token Decoding method here

	cookie, err := r.Cookie("id_token")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(cookie.Value))
	return

	token, _ := jwt.Parse(cookie.Value, func(tkn *jwt.Token) (interface{}, error){

		fmt.Println(tkn.Method)
		fmt.Println(tkn.Claims)
		fmt.Println(tkn.Signature)

		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})


	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user User
		mapstructure.Decode(claims, &user)

		json.NewEncoder(w).Encode(user)




		accRes := AccessEnforcer.Enforce("Test.Patterson@fortisureit.com", r.URL.String(), r.Method)
		if accRes != false {
			next(w, r)
		} else {
			w.WriteHeader(401)
			w.Write([]byte("User not authorized to access this resource"))
			return
		}




	} else {
		w.Write([]byte(err.Error()))
		json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
	}

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

	return result, err
}