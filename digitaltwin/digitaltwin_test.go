package digitaltwin

import (
	"azure-adt-example/azuread"
	"azure-adt-example/digitaltwin/models"
	"azure-adt-example/digitaltwin/models/rec33"
	"azure-adt-example/digitaltwin/query"
	"encoding/json"
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

func TestExecuteBuilder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/tenant1/oauth2/token" && req.Method == "POST" {
			authResponse := getValidAuthenticationResponse()
			fmt.Fprintf(w, authResponse)
		} else if strings.HasPrefix(req.RequestURI, "/query?api-version") && req.Method == "POST" {
			fmt.Fprintf(w, getSingleEntityResponseBody())
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

	builder := query.NewBuilder(rec33.Building{}, false, false)

	result, err := ExecuteBuilder[rec33.Building](client, builder)
	if err != nil {
		t.Logf("Expected nil error, but got %v", err)
		t.FailNow()
	}

	if len(result) != 2 {
		t.Logf("Expected 2 results, but got %d", len(result))
		t.FailNow()
	}

	if !find(result, func(b rec33.Building) bool { return b.ExternalId == "building01" }) {
		t.Errorf("Results does not contain expected 'building01' twin id")
	}

	if !find(result, func(b rec33.Building) bool { return b.ExternalId == "building02" }) {
		t.Errorf("Results does not contain expected 'building02' twin id")
	}
}

func TestExecuteBuilder2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/tenant1/oauth2/token" && req.Method == "POST" {
			authResponse := getValidAuthenticationResponse()
			fmt.Fprintf(w, authResponse)
		} else if strings.HasPrefix(req.RequestURI, "/query?api-version") && req.Method == "POST" {
			fmt.Fprintf(w, get2EntityResponseBody())
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

	builder := query.NewBuilder(rec33.Company{}, false, false)
	_ = builder.AddJoin(rec33.Company{}, rec33.Building{}, "owns", false, false)

	result, err := ExecuteBuilder2[rec33.Company, rec33.Building](client, builder)
	if err != nil {
		t.Logf("Expected nil error, but got %v", err)
		t.FailNow()
	}

	if len(result) != 2 {
		t.Logf("Expected 2 results, but got %d", len(result))
		t.FailNow()
	}

	if result[0].Twin1.ExternalId != "company01" {
		t.Errorf("Expected 'company01', but got '%s'", result[0].Twin1.ExternalId)
	}

	if result[1].Twin2.ExternalId != "building02" {
		t.Errorf("Expected 'building02', but got '%s'", result[1].Twin2.ExternalId)
	}
}

func TestExecuteBuilder3(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/tenant1/oauth2/token" && req.Method == "POST" {
			authResponse := getValidAuthenticationResponse()
			fmt.Fprintf(w, authResponse)
		} else if strings.HasPrefix(req.RequestURI, "/query?api-version") && req.Method == "POST" {
			fmt.Fprintf(w, get3EntityResponseBody())
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

	builder := query.NewBuilder(rec33.Company{}, false, false)
	_ = builder.AddJoin(rec33.Company{}, rec33.Building{}, "owns", false, false)
	_ = builder.AddJoin(rec33.Building{}, rec33.Level{}, "isPartOf", false, false)

	result, err := ExecuteBuilder3[rec33.Company, rec33.Building, rec33.Level](client, builder)
	if err != nil {
		t.Logf("Expected nil error, but got %v", err)
		t.FailNow()
	}

	if len(result) != 2 {
		t.Logf("Expected 2 results, but got %d", len(result))
		t.FailNow()
	}

	if result[0].Twin3.ExternalId != "level01" {
		t.Errorf("Expected 'level01', but got '%s'", result[0].Twin3.ExternalId)
	}

	if result[1].Twin3.ExternalId != "level02" {
		t.Errorf("Expected 'level02', but got '%s'", result[1].Twin3.ExternalId)
	}
}

func find[T models.IModel](list []T, f func(T) bool) bool {
	for _, v := range list {
		if f(v) {
			return true
		}
	}
	return false
}

func getValidAuthenticationResponse() string {
	return fmt.Sprintf("{ \"token_type\": \"Bearer\", \"expires_in\": \"3599\", \"ext_expires_in\": \"3599\", \"expires_on\": \"%[1]d\", \"not_before\": \"%[1]d\", \"resource\": \"0b07f429-9f4b-4714-9392-cc5e8e80c8b\", \"access_token\": \"abc123\" }", time.Now().Unix())
}

func getSingleEntityResponseBody() string {
	building1 := rec33.Building{
		GenericModel: models.GenericModel{
			ExternalId: "building01",
			ETag:       "abc123etag",
			Metadata: map[string]interface{}{
				"model": rec33.Building{}.Model(),
				"name": map[string]interface{}{
					"lastUpdated": "2022-06-22T09:09:17",
				},
			},
		},
		Name: "Test Building 1",
	}

	building2 := rec33.Building{
		GenericModel: models.GenericModel{
			ExternalId: "building02",
			ETag:       "abc456etag",
			Metadata: map[string]interface{}{
				"model": rec33.Building{}.Model(),
				"name": map[string]interface{}{
					"lastUpdated": "2021-10-21T18:09:17",
				},
			},
		},
		Name: "Test Building 2",
	}

	building1Body, _ := json.Marshal(building1)
	building2Body, _ := json.Marshal(building2)

	queryResult := QueryResultGeneric{
		Results: digitalTwinResults{
			{
				"building": building1Body,
			},
			{
				"building": building2Body,
			},
		},
		ContinuationToken: "",
	}

	resultBody, _ := json.Marshal(queryResult)

	return string(resultBody)
}

func get2EntityResponseBody() string {
	metadata := map[string]interface{}{
		"model": rec33.Building{}.Model(),
		"name": map[string]interface{}{
			"lastUpdated": "2022-06-22T09:09:17",
		},
	}

	company1 := rec33.Company{
		GenericModel: models.GenericModel{
			ExternalId: "company01",
			ETag:       "abc123etag",
			Metadata:   metadata,
		},
		Logo: "logo.svg",
		Name: "Test Company 1",
	}

	company2 := rec33.Company{
		GenericModel: models.GenericModel{
			ExternalId: "company02",
			ETag:       "abc456etag",
			Metadata:   metadata,
		},
		Logo: "logo.svg",
		Name: "Test Company 2",
	}

	building1 := rec33.Building{
		GenericModel: models.GenericModel{
			ExternalId: "building01",
			ETag:       "abc123etag",
			Metadata:   metadata,
		},
		Name: "Test Building 1",
	}

	building2 := rec33.Building{
		GenericModel: models.GenericModel{
			ExternalId: "building02",
			ETag:       "abc456etag",
			Metadata:   metadata,
		},
		Name: "Test Building 2",
	}

	company1Body, _ := json.Marshal(company1)
	company2Body, _ := json.Marshal(company2)
	building1Body, _ := json.Marshal(building1)
	building2Body, _ := json.Marshal(building2)

	queryResult := QueryResultGeneric{
		Results: digitalTwinResults{
			{
				"company":  company1Body,
				"building": building1Body,
			},
			{
				"company":  company2Body,
				"building": building2Body,
			},
		},
		ContinuationToken: "",
	}

	resultBody, _ := json.Marshal(queryResult)

	return string(resultBody)
}

func get3EntityResponseBody() string {
	metadata := map[string]interface{}{
		"model": rec33.Building{}.Model(),
		"name": map[string]interface{}{
			"lastUpdated": "2022-06-22T09:09:17",
		},
	}

	company1 := rec33.Company{
		GenericModel: models.GenericModel{
			ExternalId: "company01",
			ETag:       "abc123etag",
			Metadata:   metadata,
		},
		Logo: "logo.svg",
		Name: "Test Company 1",
	}

	company2 := rec33.Company{
		GenericModel: models.GenericModel{
			ExternalId: "company02",
			ETag:       "abc456etag",
			Metadata:   metadata,
		},
		Logo: "logo.svg",
		Name: "Test Company 2",
	}

	building1 := rec33.Building{
		GenericModel: models.GenericModel{
			ExternalId: "building01",
			ETag:       "abc123etag",
			Metadata:   metadata,
		},
		Name: "Test Building 1",
	}

	building2 := rec33.Building{
		GenericModel: models.GenericModel{
			ExternalId: "building02",
			ETag:       "abc456etag",
			Metadata:   metadata,
		},
		Name: "Test Building 2",
	}

	level1 := rec33.Level{
		GenericModel: models.GenericModel{
			ExternalId: "level01",
			ETag:       "abc123etag",
			Metadata:   metadata,
		},
		Name:            "Level 1",
		Number:          1,
		PersonCapacity:  20,
		PersonOccupancy: 10,
	}

	level2 := rec33.Level{
		GenericModel: models.GenericModel{
			ExternalId: "level02",
			ETag:       "abc456etag",
			Metadata:   metadata,
		},
		Name:            "Level 2",
		Number:          2,
		PersonCapacity:  5,
		PersonOccupancy: 8,
	}

	company1Body, _ := json.Marshal(company1)
	company2Body, _ := json.Marshal(company2)
	building1Body, _ := json.Marshal(building1)
	building2Body, _ := json.Marshal(building2)
	level1Body, _ := json.Marshal(level1)
	level2Body, _ := json.Marshal(level2)

	queryResult := QueryResultGeneric{
		Results: digitalTwinResults{
			{
				"company":  company1Body,
				"building": building1Body,
				"level":    level1Body,
			},
			{
				"company":  company2Body,
				"building": building2Body,
				"level":    level2Body,
			},
		},
		ContinuationToken: "",
	}

	resultBody, _ := json.Marshal(queryResult)

	return string(resultBody)
}
