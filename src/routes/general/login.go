package general

import (
	"math/rand"
	"net/http"
	"github.com/therecluse26/fortisure-api/src/config/mainconf"
	"strconv"
	"time"
)

var conf = mainconf.BuildConfig()

func Login(w http.ResponseWriter, r *http.Request){

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	rndInt := strconv.Itoa(r1.Intn(10000000000))

	url := conf.AuthHost + "/v2.0/authorize?client_id=" + conf.AuthClientId + "&response_type=id_token&redirect_uri=http%3A%2F%2Flocalhost:9988%2Ftoken%2F&scope=openid&state="+rndInt+"&nonce=67890&response_mode=fragment"

	w.Write([]byte(url))

}