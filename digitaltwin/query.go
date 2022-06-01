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

func queryTwin(client *Client, query string) (*[]byte, error) {
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

	return &respContent, nil
}

func QueryTwin[T models.IModel](client *Client, query string) ([]T, error) {
	queryData, err := queryTwin(client, query)
	if err != nil {
		return nil, err
	}

	var data QueryResult[T]
	err = json.Unmarshal(*queryData, &data)
	if err != nil {
		return nil, fmt.Errorf("unable to extract digital twin results: %v", err)
	}

	log.Printf("Retrieved %d records", len(data.Results))

	results := make([]T, len(data.Results))
	entityAlias := (*new(T)).Alias()
	for i, v := range data.Results {
		entry, ok := v[entityAlias]
		if ok {
			results[i] = entry
		}
	}

	return results, nil
}

func ExecuteBuilder[T1, T2 models.IModel](client *Client, builder *Builder) ([][]models.IModel, error) {
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

	results := make([][]models.IModel, len(data.Results))

	for i, v := range data.Results {
		row := make([]models.IModel, 2)

		content, ok := v[type1.Alias()]
		if ok {
			var t1 T1
			err = json.Unmarshal(content, &t1)
			row[0] = t1
		}

		if content, ok := v[type2.Alias()]; ok {
			var t2 T2
			err = json.Unmarshal(content, &t2)
			row[1] = t2
		}

		results[i] = row
	}

	log.Printf("found %d records", len(data.Results))

	return results, nil
}

func GetTwinsOfType[T models.IModel](client *Client) ([]T, error) {
	aliasName := (*new(T)).Alias()
	query := fmt.Sprintf("SELECT %[1]s FROM digitaltwins %[1]s WHERE IS_OF_MODEL(%[1]s, '%[2]s')", aliasName, (*new(T)).Model())
	return QueryTwin[T](client, query)
}
