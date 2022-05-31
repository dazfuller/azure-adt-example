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
	"strings"
	"time"
)

func queryTwin[T models.IModel](client *Client, query string) ([]T, error) {
	currentTime := time.Now().Unix()
	var err error
	endpoint := client.getQueryEndpoint()

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

	var data QueryResult[T]
	err = json.Unmarshal(respContent, &data)
	if err != nil {
		return nil, fmt.Errorf("unable to extract digital twin results: %v", err)
	}

	log.Printf("Retrieved %d records", len(data.Results))

	results := make([]T, len(data.Results))
	for i, v := range data.Results {
		entry, ok := v[GetModelAlias[T]()]
		if ok {
			results[i] = entry
		}
	}

	return results, nil
}

func GetTwinsOfType[T models.IModel](client *Client) ([]T, error) {
	aliasName := GetModelAlias[T]()
	query := fmt.Sprintf("SELECT %[1]s FROM digitaltwins %[1]s WHERE IS_OF_MODEL(%[1]s, '%[2]s')", aliasName, (*new(T)).Model())
	return queryTwin[T](client, query)
}

func GetModelAlias[T models.IModel]() string {
	var entity T
	entityNameParts := strings.Split(fmt.Sprintf("%T", entity), ".")
	entityName := strings.ToLower(entityNameParts[len(entityNameParts)-1])
	return entityName
}
