package query

import (
	"azure-adt-example/digitaltwin/models"
	"fmt"
	"strings"
)

// Builder defines a type for generating Azure Digital Twin SQL queries based on models.IModel
// types.
type Builder struct {
	from         models.IModel
	validateFrom bool
	join         []Join
	where        []IWhere
	project      []models.IModel
}

// Join represents a join condition, defining the twin being joined from and to, it's
// relationship and if the target type requires model verification.
type Join struct {
	source       models.IModel
	target       models.IModel
	relationship string
	validateType bool
}

// NewBuilder creates a new Builder type based on a required source twin and sets if that twin
// requires model type verification.
func NewBuilder(from models.IModel, validateType bool) *Builder {
	return &Builder{
		from:         from,
		validateFrom: validateType,
		join:         make([]Join, 0),
		where:        make([]IWhere, 0),
		project:      make([]models.IModel, 0),
	}
}

// AddJoin adds a new join condition to the Builder. Joins can only be specified once.
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

// WhereId defines a simple where condition that filters results to only those where a
// twin's id matches the given value.
func (b *Builder) WhereId(source models.IModel, id string) error {
	if !b.sourceExists(source) {
		return fmt.Errorf("source %s is not part of the query", source.Alias())
	}

	condition, err := NewWhereCondition(source, "ExternalId", Equals, id)
	if err != nil {
		return err
	}

	b.where = append(b.where, condition)

	return nil
}

// WhereStringFunction applies a where condition to the query
func (b *Builder) WhereStringFunction(source models.IModel, property string, function StringFunction, value string) error {
	if !b.sourceExists(source) {
		return fmt.Errorf("source %s is not part of the query", source.Alias())
	}

	whereFunction, err := NewWhereFunction[StringFunction](source, property, function, value)
	if err != nil {
		return err
	}

	b.where = append(b.where, whereFunction)

	return nil
}

func (b *Builder) WhereBooleanFunction(source models.IModel, property string, function BooleanExpressionFunction, value any) error {
	if !b.sourceExists(source) {
		return fmt.Errorf("source %s is not part of the query", source.Alias())
	}

	whereFunction, err := NewWhereFunction(source, property, function, value)
	if err != nil {
		return err
	}

	b.where = append(b.where, whereFunction)

	return nil
}

// AddProjection adds an output models.IModel type to the query, this is the equivalent of
// writing "SELECT <model type>" in the query.
func (b *Builder) AddProjection(source models.IModel) error {
	if !b.sourceExists(source) {
		return fmt.Errorf("source %s is not part of the query", source.Alias())
	}

	if !b.projectionExists(source) {
		b.project = append(b.project, source)
	}

	return nil
}

// sourceExists checks to see if a source has already been added to the builder.
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

// projectionExists checks to see if a models.IModel has already been added to
// the builder.
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

// CreateQuery takes the properties assigned to the Builder and generates a valid
// Azure Digital Twin SQL query.
func (b *Builder) CreateQuery() (string, error) {
	selectTwins := make([]string, len(b.project))
	if len(b.project) == 0 {
		selectTwins = append(selectTwins, b.from.Alias())
	} else {
		for i, p := range b.project {
			selectTwins[i] = p.Alias()
		}
	}

	whereStatements := make([]string, len(b.where))
	for i, ws := range b.where {
		whereStatements[i] = ws.GenerateClause()
	}

	fromStatement := fmt.Sprintf("digitaltwins %s", b.from.Alias())
	if b.validateFrom {
		whereStatements = append(whereStatements, ModelValidationClause(b.from).GenerateClause())
	}

	joinStatements := make([]string, len(b.join))
	for i, j := range b.join {
		joinStatements[i] = fmt.Sprintf("JOIN %s RELATED %s.%s", j.target.Alias(), j.source.Alias(), j.relationship)
		if j.validateType {
			whereStatements = append(whereStatements, ModelValidationClause(j.target).GenerateClause())
		}
	}

	joinStatement := strings.Join(joinStatements, " ")

	var whereStatement string
	if len(whereStatements) > 0 {
		whereStatement = fmt.Sprintf("WHERE %s", strings.Join(whereStatements, " AND "))
	}

	return strings.TrimSpace(fmt.Sprintf("SELECT %s FROM %s %s %s", strings.Join(selectTwins, ", "), fromStatement, joinStatement, whereStatement)), nil
}
