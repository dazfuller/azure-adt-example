package digitaltwin

import (
	"azure-adt-example/azuread"
	"azure-adt-example/digitaltwin/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// queryTwin contains the logic for querying the Azure Digital Twin instance. It returns a
// byte array of data retrieved from the API.
func queryTwin(client *Client, query string) (*[]byte, error) {
	currentTime := time.Now().Unix()
	var err error
	endpoint := client.getQueryEndpoint()

	// If the access token has not been provided, or has expired, then refresh the
	// access token
	if client.accessToken == nil || client.accessToken.ExpiresOn <= currentTime {
		client.accessToken, err = azuread.GetBearerToken(client.configuration)
		if err != nil {
			return nil, err
		}
	}

	jsonData := []byte(fmt.Sprintf(`{ "query": "%s" }`, query))

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.accessToken.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("max-items-per-page", "100")

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

// ExecuteBuilder2 queries the Azure Digital Twin using the query creating from the Builder instance. It
// returns an array of TwinResult2 objects which are typed to models.IModel types.
func ExecuteBuilder2[T1, T2 models.IModel](client *Client, builder *Builder) ([]TwinResult2[T1, T2], error) {
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

	query, err := builder.CreateQuery()
	if err != nil {
		return nil, err
	}

	queryData, err := queryTwin(client, query)
	if err != nil {
		return nil, err
	}

	var data QueryResultGeneric
	err = json.Unmarshal(*queryData, &data)
	if err != nil {
		return nil, fmt.Errorf("unable to extract digital twin results: %v", err)
	}

	results := make([]TwinResult2[T1, T2], len(data.Results))

	for i, v := range data.Results {
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

	log.Printf("found %d records", len(data.Results))

	return results, nil
}
