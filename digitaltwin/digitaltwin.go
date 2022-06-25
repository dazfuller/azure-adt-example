package digitaltwin

import (
	"azure-adt-example/azuread"
	"azure-adt-example/digitaltwin/models"
	"azure-adt-example/digitaltwin/query"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const apiVersion = "2020-10-31"

// Client defines the values required for accessing an Azure Digital Twin instance.
type Client struct {
	configuration   *azuread.TwinConfiguration
	accessToken     *azuread.AccessToken
	MaxItemsPerPage uint
}

// NewClient creates an instance of the Client type.
func NewClient(configuration *azuread.TwinConfiguration, accessToken *azuread.AccessToken) *Client {
	client := Client{configuration: configuration, accessToken: accessToken, MaxItemsPerPage: 1000}
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

func (c *Client) getBuilderResults(builder *query.Builder, err error) (digitalTwinResults, error) {
	queryResults := make(digitalTwinResults, 0)
	var continuationToken *string
	generatedQuery, err := builder.CreateQuery()
	if err != nil {
		return nil, fmt.Errorf("unable to generate digital twin query: %s", err)
	}

	for {
		queryData, err := c.queryTwin(*generatedQuery, continuationToken)
		if err != nil {
			return nil, err
		}

		var data QueryResultGeneric
		err = json.Unmarshal(*queryData, &data)
		if err != nil {
			return nil, fmt.Errorf("unable to extract digital twin results: %v", err)
		}

		queryResults = append(queryResults, data.Results...)

		if !data.HasContinuationToken() {
			break
		}

		continuationToken = &data.ContinuationToken
	}
	log.Printf("Total number of records: %d", len(queryResults))
	return queryResults, nil
}

// queryTwin contains the logic for querying the Azure Digital Twin instance. It returns a
// byte array of data retrieved from the API.
func (c *Client) queryTwin(query string, continuationToken *string) (*[]byte, error) {
	currentTime := time.Now().Unix()
	var err error
	endpoint := c.getQueryEndpoint()

	// If the access token has not been provided, or has expired, then refresh the
	// access token
	if c.accessToken == nil || c.accessToken.ExpiresOn <= currentTime {
		c.accessToken, err = azuread.GetBearerToken(c.configuration)
		if err != nil {
			return nil, err
		}
	}

	var requestBody string
	if continuationToken == nil {
		requestBody = fmt.Sprintf(`{ "query": "%s" }`, query)
	} else {
		continuationData, _ := json.Marshal(*continuationToken)
		requestBody = fmt.Sprintf(`{ "query": "%s", "continuationToken": %s }`, query, string(continuationData))
	}

	jsonData := []byte(requestBody)
	maxItemsPerPage := c.MaxItemsPerPage
	if maxItemsPerPage == 0 {
		maxItemsPerPage = 1
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("max-items-per-page", fmt.Sprint(maxItemsPerPage))

	log.Printf("Querying endpoint %s with query:\n%s", endpoint, query)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		respContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("non-success status code returned: %d", resp.StatusCode)
		} else {
			var respError QueryError
			_ = json.Unmarshal(respContent, &respError)
			return nil, fmt.Errorf("non-success status code returned: %d\nQuery: %s\n%s", resp.StatusCode, query, respError.ErrorDetail.Message)
		}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Unable to close body: %v", err)
		}
	}(resp.Body)

	log.Printf("Response code: %d", resp.StatusCode)

	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &respContent, nil
}

// ExecuteBuilder queries the Azure Digital Twin using the query creating from the Builder instance. It
// returns an array of models.IModel types.
func ExecuteBuilder[T1 models.IModel](client *Client, builder *query.Builder) ([]T1, error) {
	type1 := *new(T1)

	var err error

	if err = builder.AddProjection(type1); err != nil {
		return nil, err
	}

	queryResults, err := client.getBuilderResults(builder, err)
	if err != nil {
		return nil, err
	}

	results := make([]T1, len(queryResults))

	for i, v := range queryResults {
		t1 := new(T1)

		content, ok := v[type1.Alias()]
		if ok {
			err = json.Unmarshal(content, t1)
			if err != nil {
				return nil, fmt.Errorf("unable to parse %v into %T", content, *t1)
			}
			results[i] = *t1
		}
	}

	return results, nil
}

// ExecuteBuilder2 queries the Azure Digital Twin using the query creating from the Builder instance. It
// returns an array of TwinResult2 objects which are typed to models.IModel types.
func ExecuteBuilder2[T1, T2 models.IModel](client *Client, builder *query.Builder) ([]TwinResult2[T1, T2], error) {
	type1 := *new(T1)
	type2 := *new(T2)

	err := builder.AddProjection(type1)
	if err != nil {
		return nil, err
	}

	err = builder.AddProjection(type2)
	if err != nil {
		return nil, err
	}

	queryResults, err := client.getBuilderResults(builder, err)
	if err != nil {
		return nil, err
	}

	results := make([]TwinResult2[T1, T2], len(queryResults))

	for i, v := range queryResults {
		t1 := new(T1)
		t2 := new(T2)

		content, ok := v[type1.Alias()]
		if ok {
			err = json.Unmarshal(content, t1)
			if err != nil {
				return nil, fmt.Errorf("unable to parse %v into %T", content, *t1)
			}
		}

		if content, ok := v[type2.Alias()]; ok {
			err = json.Unmarshal(content, t2)
			if err != nil {
				return nil, fmt.Errorf("unable to parse %v into %T", content, *t2)
			}
		}

		results[i] = NewTwinResult2(t1, t2)
	}

	return results, nil
}

// ExecuteBuilder3 queries the Azure Digital Twin using the query creating from the Builder instance. It
// returns an array of TwinResult3 objects which are typed to models.IModel types.
func ExecuteBuilder3[T1, T2, T3 models.IModel](client *Client, builder *query.Builder) ([]TwinResult3[T1, T2, T3], error) {
	type1 := *new(T1)
	type2 := *new(T2)
	type3 := *new(T3)

	err := builder.AddProjection(type1)
	if err != nil {
		return nil, err
	}

	err = builder.AddProjection(type2)
	if err != nil {
		return nil, err
	}

	err = builder.AddProjection(type3)
	if err != nil {
		return nil, err
	}

	queryResults, err := client.getBuilderResults(builder, err)
	if err != nil {
		return nil, err
	}

	results := make([]TwinResult3[T1, T2, T3], len(queryResults))

	for i, v := range queryResults {
		t1 := new(T1)
		t2 := new(T2)
		t3 := new(T3)

		content, ok := v[type1.Alias()]
		if ok {
			err = json.Unmarshal(content, t1)
			if err != nil {
				return nil, fmt.Errorf("unable to parse %v into %T", content, *t1)
			}
		}

		if content, ok := v[type2.Alias()]; ok {
			err = json.Unmarshal(content, t2)
			if err != nil {
				return nil, fmt.Errorf("unable to parse %v into %T", content, *t2)
			}
		}

		if content, ok := v[type3.Alias()]; ok {
			err = json.Unmarshal(content, t3)
			if err != nil {
				return nil, fmt.Errorf("unable to parse %v into %T", content, *t3)
			}
		}

		results[i] = NewTwinResult3(t1, t2, t3)
	}

	return results, nil
}
