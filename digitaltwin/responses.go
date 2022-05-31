package digitaltwin

import "azure-adt-example/digitaltwin/models"

type QueryError struct {
	ErrorDetail ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type QueryResult[T models.IModel] struct {
	Results           []map[string]T `json:"value"`
	ContinuationToken string         `json:"continuationToken"`
}
