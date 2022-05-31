package azuread

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type AccessToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in,string"`
	ExtExpiresIn int64  `json:"ext_expires_in,string"`
	ExpiresOn    int64  `json:"expires_on,string"`
	NotBefore    int64  `json:"not_before,string"`
	Resource     string `json:"resource"`
	AccessToken  string `json:"access_token"`
}

func GetBearerToken(configuration *TwinConfiguration) (*AccessToken, error) {
	log.Printf("Attempting to acquire access token for resource: %s", configuration.ResourceId)

	authenticationUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", configuration.TenantId)

	data := url.Values{}
	data.Add("client_id", configuration.ClientId)
	data.Add("client_secret", configuration.ClientSecret)
	data.Add("resource", configuration.ResourceId)
	data.Add("grant_type", "client_credentials")

	client := &http.Client{}
	resp, err := client.PostForm(authenticationUrl, data)
	if err != nil {
		return nil, fmt.Errorf("unable to obtain access token: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %v", err)
	}

	var response AccessToken
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("unable to parse response: %v", err)
	}

	log.Print("Successfully acquired access token")

	return &response, nil
}
