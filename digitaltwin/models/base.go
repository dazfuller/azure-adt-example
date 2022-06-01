package models

import (
	"fmt"
	"reflect"
	"strings"
)

type IModel interface {
	Model() string
	Alias() string
	ValidationClause() string
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

func ModelValidationClause[T IModel]() string {
	t := *new(T)
	return fmt.Sprintf("IS_OF_MODEL(%s, '%s')", t.Alias(), t.Model())
}
