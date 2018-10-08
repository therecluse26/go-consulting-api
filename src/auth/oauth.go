package auth

import (
	"golang.org/x/oauth2"
	"fmt"
	"context"
	"golang.org/x/oauth2/microsoft"
)

func OAuth(redir_url string, code string) {

	tenant := "b511c547-996d-4ea5-a743-64f4486b22bc"
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     "c7ba4700-1b55-4563-b066-9d103d59efcc",
		ClientSecret: "ZhBrzmAXreycvk+9s0Eqx58FPF0wO3a5Y2TiccTIgms=",
		Scopes:       []string{"openid"},
		Endpoint: microsoft.AzureADEndpoint(tenant),
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	/*url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	fmt.Printf("Visit the URL for the auth dialog: %v", url)*/

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.

	/*var code string
	if _, err := fmt.Scan(&code); err != nil {
		fmt.Println(err)
	}*/

	fmt.Println("Code: ")

	fmt.Println(code)

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(tok)

	client := conf.Client(ctx, tok)
	client.Get(redir_url)

}