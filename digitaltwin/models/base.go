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

func (gm *GenericModel) TwinModelType() string {
	model, ok := gm.Metadata["$model"]
	if ok {
		return fmt.Sprintf("%s", model)
	}

	return "Unknown"
}

func GetModelAlias[T IModel]() string {
	typeNameParts := strings.Split(reflect.TypeOf(*new(T)).Name(), ",")
	return strings.ToLower(typeNameParts[len(typeNameParts)-1])
}
