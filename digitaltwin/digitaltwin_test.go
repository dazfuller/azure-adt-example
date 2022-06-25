package digitaltwin

import (
	"azure-adt-example/azuread"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestClient_getQueryEndpoint(t *testing.T) {
	testUrl, _ := url.Parse("https://myenv.reg.digitaltwins.azure.net")

	conf := azuread.TwinConfiguration{
		URL:          *testUrl,
		ClientId:     "client_id",
		ClientSecret: "client_secret",
		TenantId:     "tenant_id",
		ResourceId:   "1234",
	}

	token := azuread.AccessToken{AccessToken: "abc123"}

	client := NewClient(&conf, &token)

	expected := "https://myenv.reg.digitaltwins.azure.net/query?api-version=2020-10-31"

	actual := client.getQueryEndpoint()

	if actual != expected {
		t.Errorf("Expected endpoint '%s', but got '%s'", expected, actual)
	}
}

func TestClient_queryTwin(t *testing.T) {
	expectedQueryResponse := "Expected response"
	expectedBody := "{ \"query\": \"SELECT * FROM digitaltwins\" }"
	var queryRequest *http.Request
	var queryBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/tenant1/oauth2/token" && req.Method == "POST" {
			authResponse := getValidAuthenticationResponse()
			fmt.Fprintf(w, authResponse)
		} else if strings.HasPrefix(req.RequestURI, "/query?api-version") && req.Method == "POST" {
			fmt.Fprintf(w, expectedQueryResponse)
			queryRequest = req
			bodyData, _ := ioutil.ReadAll(req.Body)
			queryBody = string(bodyData)
		}
	}))
	defer server.Close()

	serverUrl, _ := url.Parse(server.URL)

	conf := azuread.TwinConfiguration{
		URL:          *serverUrl,
		ClientId:     "client1",
		ClientSecret: "secret1",
		TenantId:     "tenant1",
		ResourceId:   "resource1",
		AuthorityUrl: *serverUrl,
	}

	token := azuread.AccessToken{AccessToken: "abc123"}

	client := NewClient(&conf, &token)

	data, err := client.queryTwin("SELECT * FROM digitaltwins", nil)
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
		t.FailNow()
	}

	if !reflect.DeepEqual(*data, []byte(expectedQueryResponse)) {
		t.Errorf("Expected response '%s', but got '%s'", expectedQueryResponse, string(*data))
		t.FailNow()
	}

	if queryRequest.Header.Get("Authorization") != "Bearer abc123" {
		t.Errorf("Expected Authorization header of 'Bearer abc123', but got '%s'", queryRequest.Header.Get("Authorization"))
	}

	if queryRequest.Header.Get("Max-Items-Per-Page") != "1000" {
		t.Errorf("Expected 100 items per page header value, but got %s", queryRequest.Header.Get("Max-Items-Per-Page"))
	}

	if queryBody != expectedBody {
		t.Errorf("Expected to get body '%s', but got '%s'", expectedBody, queryBody)
	}
}

func TestClient_queryTwin_failedAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/tenant1/oauth2/token" && req.Method == "POST" {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Invalid authentication request")
		} else if strings.HasPrefix(req.RequestURI, "/query?api-version") && req.Method == "POST" {
			fmt.Fprintf(w, "Should not get here")
		}
	}))
	defer server.Close()

	serverUrl, _ := url.Parse(server.URL)

	conf := azuread.TwinConfiguration{
		URL:          *serverUrl,
		ClientId:     "client1",
		ClientSecret: "secret1",
		TenantId:     "tenant1",
		ResourceId:   "resource1",
		AuthorityUrl: *serverUrl,
	}

	client := NewClient(&conf, nil)

	expectedError := "received error response from authority url: 400"

	_, err := client.queryTwin("SELECT * FROM digitaltwins", nil)
	if err == nil {
		t.Error("Expected error, but got nil")
	} else if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error '%s', but got '%s", expectedError, err)
	}
}

func TestClient_queryTwin_ContinuationToken(t *testing.T) {
	var queryBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/tenant1/oauth2/token" && req.Method == "POST" {
			authResponse := getValidAuthenticationResponse()
			fmt.Fprintf(w, authResponse)
		} else if strings.HasPrefix(req.RequestURI, "/query?api-version") && req.Method == "POST" {
			fmt.Fprintf(w, "Expected response")
			bodyData, _ := ioutil.ReadAll(req.Body)
			queryBody = string(bodyData)
		}
	}))
	defer server.Close()

	serverUrl, _ := url.Parse(server.URL)

	conf := azuread.TwinConfiguration{
		URL:          *serverUrl,
		ClientId:     "client1",
		ClientSecret: "secret1",
		TenantId:     "tenant1",
		ResourceId:   "resource1",
		AuthorityUrl: *serverUrl,
	}

	token := azuread.AccessToken{AccessToken: "abc123"}
	query := "SELECT * FROM digitaltwins"
	continuationToken := "{ \"continuationToken\": \"some token value\" }"

	expectedBody := fmt.Sprintf("{ \"query\": \"%s\", \"continuationToken\": \"{ \\\"continuationToken\\\": \\\"some token value\\\" }\" }", query)

	client := NewClient(&conf, &token)

	_, err := client.queryTwin(query, &continuationToken)
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
		t.FailNow()
	}

	if queryBody != expectedBody {
		t.Errorf("Expected body '%s', but got '%s'", expectedBody, queryBody)
	}
}

func TestClient_queryTwin_FailedRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/tenant1/oauth2/token" && req.Method == "POST" {
			authResponse := getValidAuthenticationResponse()
			fmt.Fprintf(w, authResponse)
		} else if strings.HasPrefix(req.RequestURI, "/query?api-version") && req.Method == "POST" {
			w.WriteHeader(404)
			w.Write([]byte("{ \"error\": { \"code\": \"InvalidTwin\", \"message\": \"The request twin was not found\" } }"))
		}
	}))
	defer server.Close()

	serverUrl, _ := url.Parse(server.URL)

	conf := azuread.TwinConfiguration{
		URL:          *serverUrl,
		ClientId:     "client1",
		ClientSecret: "secret1",
		TenantId:     "tenant1",
		ResourceId:   "resource1",
		AuthorityUrl: *serverUrl,
	}

	token := azuread.AccessToken{AccessToken: "abc123"}
	query := "SELECT * FROM digitaltwins"

	expectedError := "non-success status code returned: 404\nQuery: SELECT * FROM digitaltwins\nThe request twin was not found"

	client := NewClient(&conf, &token)

	_, err := client.queryTwin(query, nil)
	if err == nil {
		t.Logf("Expected error but got nil")
		t.FailNow()
	}

	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', but got '%s'", expectedError, err.Error())
	}
}

func getValidAuthenticationResponse() string {
	return fmt.Sprintf("{ \"token_type\": \"Bearer\", \"expires_in\": \"3599\", \"ext_expires_in\": \"3599\", \"expires_on\": \"%[1]d\", \"not_before\": \"%[1]d\", \"resource\": \"0b07f429-9f4b-4714-9392-cc5e8e80c8b\", \"access_token\": \"abc123\" }", time.Now().Unix())
}
