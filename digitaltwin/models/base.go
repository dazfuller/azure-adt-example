package models

import (
	"fmt"
	"reflect"
	"strings"
)

type IModel interface {
	Model() string
	Alias() string
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

func (gm GenericModel) Alias() string {
	typeNameParts := strings.Split(reflect.TypeOf(gm).Name(), ",")
	return typeNameParts[len(typeNameParts)-1]
}
