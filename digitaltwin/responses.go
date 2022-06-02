package digitaltwin

import (
	"azure-adt-example/digitaltwin/models"
	"encoding/json"
)

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

type QueryResultGeneric struct {
	Results           []map[string]json.RawMessage `json:"value"`
	ContinuationToken string                       `json:"continuationToken"`
}

type TwinResult2[T1, T2 models.IModel] struct {
	Twin1 T1
	Twin2 T2
}

func NewTwinResult2[T1, T2 models.IModel](t1 *T1, t2 *T2) TwinResult2[T1, T2] {
	return TwinResult2[T1, T2]{*t1, *t2}
}
