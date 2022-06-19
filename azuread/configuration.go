package azuread

import (
	"github.com/subosito/gotenv"
	"log"
	"net/url"
	"os"
)

const (
	ResourceId   = "0b07f429-9f4b-4714-9392-cc5e8e80c8b0"
	AuthorityUrl = "https://login.microsoftonline.com"
)

// TwinConfiguration defines properties required for connecting to an Azure Digital
// Twin instance.
type TwinConfiguration struct {
	// URL defines the endpoint of the twin instance to connect to.
	URL url.URL

	// ClientId of the service principal account used for connecting to the Azure
	// Digital Twin.
	ClientId string

	// ClientSecret of the service principal account used for connecting to the
	// Azure Digital Twin.
	ClientSecret string

	// TenantId (Directory ID) where the service principal authenticates to.
	TenantId string

	// ResourceId which defines the scope of the AccessToken when it's retrieved.
	ResourceId string

	// AuthorityUrl defines the base url for obtaining an access token (e.g. https://login.microsoftonline.com)
	AuthorityUrl url.URL
}

// NewTwinConfiguration creates a new instance of TwinConfiguration
func NewTwinConfiguration() *TwinConfiguration {
	err := gotenv.Load()
	if err != nil {
		log.Printf("Unable to load environment variables: %v", err)
	}

	twinUrl, err := url.Parse(getEnvironmentValue("TWIN_URL"))
	if err != nil {
		log.Fatal("Invalid Twin URL in environment variable TWIN_URL")
	}

	authority, _ := url.Parse(AuthorityUrl)

	return &TwinConfiguration{
		*twinUrl,
		getEnvironmentValue("TWIN_CLIENT_ID"),
		getEnvironmentValue("TWIN_CLIENT_SECRET"),
		getEnvironmentValue("TWIN_TENANT_ID"),
		ResourceId,
		*authority,
	}
}

func getEnvironmentValue(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s not set in environment variables", key)
	}
	return value
}
