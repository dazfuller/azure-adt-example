package digitaltwin

import (
	"azure-adt-example/digitaltwin/models"
	"fmt"
	"strings"
)

type Builder struct {
	from         models.IModel
	validateFrom bool
	join         []Join
	where        []string
	project      []models.IModel
}

type Join struct {
	source       models.IModel
	target       models.IModel
	relationship string
	validateType bool
}

func NewBuilder(from models.IModel, validateType bool) *Builder {
	return &Builder{
		from:         from,
		validateFrom: validateType,
		join:         make([]Join, 0),
		where:        make([]string, 0),
		project:      make([]models.IModel, 0),
	}
}

func (b *Builder) AddJoin(source models.IModel, target models.IModel, relationship string, validateType bool) error {
	exists := false
	for _, j := range b.join {
		if target.Alias() == j.target.Alias() {
			exists = true
			break
		}
	}

	if exists {
		return fmt.Errorf("a target of alias '%s' already exists", target.Alias())
	}

	join := Join{
		source:       source,
		target:       target,
		relationship: relationship,
		validateType: validateType,
	}

	b.join = append(b.join, join)

	return nil
}

func (b *Builder) WhereId(source models.IModel, id string) error {
	if !b.sourceExists(source) {
		return fmt.Errorf("source %s is not part of the query", source.Alias())
	}

	b.where = append(b.where, fmt.Sprintf("%s.$dtId = '%s'", source.Alias(), id))

	return nil
}

func (b *Builder) AddProjection(source models.IModel) error {
	if !b.sourceExists(source) {
		return fmt.Errorf("source %s is not part of the query", source.Alias())
	}

	if !b.projectionExists(source) {
		b.project = append(b.project, source)
	}

	return nil
}

func (b *Builder) sourceExists(source models.IModel) bool {
	sourceExists := b.from.Alias() == source.Alias()

	if !sourceExists {
		for _, j := range b.join {
			if j.target.Alias() == source.Alias() {
				sourceExists = true
				break
			}
		}
	}

	return sourceExists
}

func (b *Builder) projectionExists(source models.IModel) bool {
	exists := false

	for _, p := range b.project {
		if p.Alias() == source.Alias() {
			exists = true
			break
		}
	}

	return exists
}

func (b *Builder) CreateQuery() (string, error) {
	selectTwins := make([]string, len(b.project))
	if len(b.project) == 0 {
		selectTwins = append(selectTwins, b.from.Alias())
	} else {
		for i, p := range b.project {
			selectTwins[i] = p.Alias()
		}
	}

	whereStatements := b.where

	fromStatement := fmt.Sprintf("digitaltwins %s", b.from.Alias())
	if b.validateFrom {
		whereStatements = append(whereStatements, b.from.ValidationClause())
	}

	joinStatements := make([]string, len(b.join))
	for i, j := range b.join {
		joinStatements[i] = fmt.Sprintf("JOIN %s RELATED %s.%s", j.target.Alias(), j.source.Alias(), j.relationship)
		if j.validateType {
			whereStatements = append(whereStatements, j.target.ValidationClause())
		}
	}

	joinStatement := strings.Join(joinStatements, " ")

	var whereStatement string
	if len(whereStatements) > 0 {
		whereStatement = fmt.Sprintf("WHERE %s", strings.Join(whereStatements, " AND "))
	}

	return strings.TrimSpace(fmt.Sprintf("SELECT %s FROM %s %s %s", strings.Join(selectTwins, ", "), fromStatement, joinStatement, whereStatement)), nil
}
