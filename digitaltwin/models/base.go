package models

import "fmt"

type IModel interface {
	Model() string
}

type GenericModel struct {
	ExternalId string                 `json:"$dtId"`
	ETag       string                 `json:"$etag"`
	Metadata   map[string]interface{} `json:"$metadata"`
}

func (GenericModel) Model() string {
	return "Unknown"
}

func (gm *GenericModel) TwinModelType() string {
	model, ok := gm.Metadata["$model"]
	if ok {
		return fmt.Sprintf("%s", model)
	}

	return gm.Model()
}
