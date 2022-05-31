package digitaltwin

import (
	"azure-adt-example/azuread"
	"net/url"
)

const apiVersion = "2020-10-31"

type Client struct {
	configuration *azuread.TwinConfiguration
	accessToken   *azuread.AccessToken
}

func NewClient(configuration *azuread.TwinConfiguration, accessToken *azuread.AccessToken) *Client {
	client := Client{configuration: configuration, accessToken: accessToken}
	return &client
}

func (c *Client) getQueryEndpoint() string {
	endpoint := c.configuration.URL
	endpoint.Path = "/query"

	params := url.Values{}
	params.Add("api-version", apiVersion)

	endpoint.RawQuery = params.Encode()

	return endpoint.String()
}
