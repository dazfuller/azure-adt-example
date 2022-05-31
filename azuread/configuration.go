package azuread

import (
	"github.com/subosito/gotenv"
	"log"
	"net/url"
	"os"
)

const ResourceId = "0b07f429-9f4b-4714-9392-cc5e8e80c8b0"

type TwinConfiguration struct {
	URL          url.URL
	ClientId     string
	ClientSecret string
	TenantId     string
	ResourceId   string
}

func NewTwinConfiguration() *TwinConfiguration {
	err := gotenv.Load()
	if err != nil {
		log.Fatal("Unable to load environment variables", err)
	}

	twinUrl, err := url.Parse(getEnvironmentValue("TWIN_URL"))
	if err != nil {
		log.Fatal("Invalid Twin URL in environment variable TWIN_URL")
	}

	return &TwinConfiguration{
		*twinUrl,
		getEnvironmentValue("TWIN_CLIENT_ID"),
		getEnvironmentValue("TWIN_CLIENT_SECRET"),
		getEnvironmentValue("TWIN_TENANT_ID"),
		ResourceId,
	}
}

func getEnvironmentValue(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s not set in environment variables", key)
	}
	return value
}
