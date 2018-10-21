package auth

import (
	"../database"
	"net/http"
	"fmt"
	"../routes"
)

/**
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

func mergePathParam(urlPath string, routePath string){
	/*
	for k, v := range routePath.Params {
		sqlMerged = strings.Replace(sqlMerged, "{{" + k + "}}", v, -1)
	}
	*/
}

func ProtectedEndpoint (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	path := r.URL.Path

	objName := routes.GetObjectFromPath(path)

	accResult, err := CheckUserAccess("Brad.Magyar", objName, "GET")
	if err != nil {
		fmt.Println(err)
	}

	if accResult != 0 {

		next(w, r)

	} else {

		w.WriteHeader(401)
		w.Write([]byte("User not authorized to access this resource"))
		return
	}


	/* cookie, err := req.Cookie("auth_token")
	if err != nil {
		fmt.Println(err)
	}

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
	} else {
		json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
	}*/

}