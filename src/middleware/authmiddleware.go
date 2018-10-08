package middleware

import (
	"net/http"
	"../config/mainconf"
	"../auth"
	"fmt"
)

var AuthConf = mainconf.GetAuthConfig()

var Conf = mainconf.BuildConfig()

const (
	ErrTokenInvalid = "Token is invalid"
	ErrTokenMissing = "No token is present"
	ErrTokenExpired = "Token is expired"
	ErrAccessDenied = "Access to this resource is denied"
)

func AuthMiddleware(rw http.ResponseWriter, r *http.Request, nextFunc http.HandlerFunc) {

/* 1) Check for existing token. If not found, throw error */

	/*token, _ := auth.RefreshToken(AuthConf.AuthHost, AuthConf.AuthClientId, AuthConf.AuthSecret, "openid profile email https://graph.microsoft.com/Mail.Read https://graph.microsoft.com/User.Read", "OAQABAAAAAADXzZ3ifr-GRbDT45zNSEFE9ZsnZ9JAS_29gqVxSVN0JTstowRd9JcJWNeJslVzxLPiKvZpF2VQITjh0uvizVnu61_XAxSfOf7VdFq-90VOWz_bu1v454M2K5XHFihgjECg6OYJKNUVUkJDH7SCk_hp65FTbBWGayLRqCZTtyyokKagIFR4zHYvcHZwc8IX0p8Zhu7mWQ-5RJahCHFd3B_H7yphbRPrYyP9RwFJtwBF612zuiQP--uc5w8RqvWc6wz2VRKtCiY388Qr5ikOugQHsS08J9WRZ3_L4n9sl9pFeRp4O3uiKM4INUkR0jdmqhqgFN_j8TgXlQAdOpLtbwMF8gvnD-Ba1W99Z5Srd9bV3Lpu5VuF2DGWM1gCIxrT_4liZbsPE0lVYHRMx-3gUcsI_jiu_rrJAaMLuw5Ri8Za9i4d6oV4A_c-Dyz6qFAi1sRy4eHoRXdjo1rjJbdhyZv29-zKF0xj-CKuM5VgS6nMuWaIc28RJ6vLyZGfz-FvzTT4PxoTMLJp9ZhD_SKrb57sGtqk03TMuZ3Lw00Ou2Oi6pX85ZmNREm0U69YJ8qcD0ch8QU5nkSAlYPuN_ekTzzx0tz8RoPPEP73T6ML4RrQ1Eo5mJam19JB5FUxYZLY9nDRBUg16kBe64mtA1aukHH8ElO0npX1F0QKmG2_trRG0B6HCSKmg8mbr_2CoM0IIMmikEvrdVUBNFAsMFVz88ZK58jqPMXugj6iMhAzH7WkoGslY2n46cYID8bS6Qjw7gg8hdVF0udoLv9iS7AlH950GPrKP1Yu-5QIIu17FgJbKUmNEkp99GY55sVYdv9vBr4gAA")
	fmt.Println(token)*/

	fwdUrl := Conf.ApiHost + r.RequestURI

	code := r.URL.Query().Get("code")

	fmt.Println(code)

	auth.OAuth(fwdUrl, code)

	/* 3) Check access against application */

	/* 4) Execute request */
	//nextFunc(rw, r)

	/* 5) Issue refresh token */
	//fmt.Println("auth middleware -> after the controller was executed")

}
