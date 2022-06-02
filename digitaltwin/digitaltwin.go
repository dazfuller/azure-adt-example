package digitaltwin

import (
	"azure-adt-example/azuread"
	"net/url"
)

const apiVersion = "2020-10-31"

// Client defines the values required for accessing an Azure Digital Twin instance.
type Client struct {
	configuration *azuread.TwinConfiguration
	accessToken   *azuread.AccessToken
}

// NewClient creates an instance of the Client type.
func NewClient(configuration *azuread.TwinConfiguration, accessToken *azuread.AccessToken) *Client {
	client := Client{configuration: configuration, accessToken: accessToken}
	return &client
}

// getQueryEndpoint generates the full URL required for querying the Azure Digital Twin
// instance defined in the Client configuration.
func (c *Client) getQueryEndpoint() string {
	endpoint := c.configuration.URL
	endpoint.Path = "/query"

	params := url.Values{}
	params.Add("api-version", apiVersion)

	endpoint.RawQuery = params.Encode()

	return endpoint.String()
}
