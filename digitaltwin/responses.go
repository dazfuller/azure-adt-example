package digitaltwin

import (
	"encoding/json"
)

// QueryError defines the response received from Azure Digital Twin if a query
// results in an error being returned.
type QueryError struct {
	ErrorDetail ErrorDetail `json:"error"`
}

// ErrorDetail defines the contents of an Azure Digital Twin error response
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// QueryResultGeneric defines a successful response message from Azure Digital Twin.
// The response includes a ContinuationToken for paging results. The Results are
// left as their JSON message values so that each may be Unmarshalled to the correct
// type.
type QueryResultGeneric struct {
	// Results contains an array of JSON objects where each object represents a specific
	// models.IModel type.
	Results []map[string]json.RawMessage `json:"value"`

	// ContinuationToken is the value required to be sent back to the query API to
	// retrieve the next page of results. If this is empty then no further results
	// are available.
	ContinuationToken string `json:"continuationToken"`
}
