package influxqb

import (
	"regexp"

	"github.com/influxdata/influxql"
)

type ShowContinuousQueries struct{}

func (b *ShowContinuousQueries) Build() (string, error) {
	q := &influxql.ShowContinuousQueriesStatement{}
	return q.String(), nil
}

type ShowDatabases struct{}

func (b *ShowDatabases) Build() (string, error) {
	q := &influxql.ShowDatabasesStatement{}
	return q.String(), nil
}

type ShowUsers struct{}

func (b *ShowUsers) Build() (string, error) {
	q := &influxql.ShowUsersStatement{}
	return q.String(), nil
}

type ShowQueries struct{}

func (b *ShowQueries) Build() (string, error) {
	q := &influxql.ShowQueriesStatement{}
	return q.String(), nil
}

type ShowShardGroups struct{}

func (b *ShowShardGroups) Build() (string, error) {
	q := &influxql.ShowShardGroupsStatement{}
	return q.String(), nil
}

type ShowSubscriptions struct{}

func (b *ShowSubscriptions) Build() (string, error) {
	q := &influxql.ShowSubscriptionsStatement{}
	return q.String(), nil
}

type ShowShards struct{}

func (b *ShowShards) Build() (string, error) {
	statement := &influxql.ShowShardsStatement{}
	return statement.String(), nil
}

type ShowDiagnostics struct {
	module string
}

func (b *ShowDiagnostics) ForModule(str string) *ShowDiagnostics {
	b.module = str
	return b
}

func (b *ShowDiagnostics) Build() (string, error) {
	t := &influxql.ShowDiagnosticsStatement{Module: b.module}
	return t.String(), nil
}

type ShowFieldKeyCardinality struct {
	q *influxql.ShowFieldKeyCardinalityStatement
}

func (b *ShowFieldKeyCardinality) Exact() *ShowFieldKeyCardinality {
	b.q.Exact = true
	return b
}

func (b *ShowFieldKeyCardinality) Limit(limit int) *ShowFieldKeyCardinality {
	b.q.Limit = limit
	return b
}

func (b *ShowFieldKeyCardinality) Offset(offset int) *ShowFieldKeyCardinality {
	b.q.Offset = offset
	return b
}

func (b *ShowFieldKeyCardinality) OnDatabase(str string) *ShowFieldKeyCardinality {
	b.q.Database = str
	return b
}

func (b *ShowFieldKeyCardinality) From(sources ...*Measurement) *ShowFieldKeyCardinality {
	for _, me := range sources {
		b.q.Sources = append(b.q.Sources, me.m)
	}
	return b
}

func (b *ShowFieldKeyCardinality) GroupBy(fields ...GroupByIf) *ShowFieldKeyCardinality {
	for _, f := range fields {
		if finalField := f.groupBy(); finalField != nil {
			b.q.Dimensions = append(b.q.Dimensions, finalField)
		}
	}
	return b
}

func (b *ShowFieldKeyCardinality) Where(condition interface{}) *ShowFieldKeyCardinality {
	switch condition := condition.(type) {
	case influxql.Expr:
		b.q.Condition = &influxql.ParenExpr{Expr: condition}

	case MathExprIf:
		b.q.Condition = condition.expr()

	}
	return b
}

func (b *ShowFieldKeyCardinality) Build() (string, error) {
	return b.q.String(), nil
}

type ShowFieldKeys struct {
	q *influxql.ShowFieldKeysStatement
}

func (b *ShowFieldKeys) WithDatabse(database string) *ShowFieldKeys {
	b.q.Database = database
	return b
}

func (b *ShowFieldKeys) From(meaurements ...*Measurement) *ShowFieldKeys {
	for _, m := range meaurements {
		b.q.Sources = append(b.q.Sources, m.m)
	}
	return b
}

func (b *ShowFieldKeys) Limit(limit int) *ShowFieldKeys {
	b.q.Limit = limit
	return b
}

func (b *ShowFieldKeys) Offset(offset int) *ShowFieldKeys {
	b.q.Offset = offset
	return b
}

func (b *ShowFieldKeys) Build() (string, error) {
	return b.q.String(), nil
}

type ShowGrants struct {
	q *influxql.ShowGrantsForUserStatement
}

func (b *ShowGrants) ForUser(user string) *ShowGrants {
	b.q.Name = user
	return b
}

func (b *ShowGrants) Build() (string, error) {
	return b.q.String(), nil
}

type ShowMeasurementCardinality struct {
	q *influxql.ShowMeasurementCardinalityStatement
}

func (b *ShowMeasurementCardinality) Exact() *ShowMeasurementCardinality {
	b.q.Exact = true
	return b
}

func (b *ShowMeasurementCardinality) Limit(limit int) *ShowMeasurementCardinality {
	b.q.Limit = limit
	return b
}

func (b *ShowMeasurementCardinality) Offset(offset int) *ShowMeasurementCardinality {
	b.q.Offset = offset
	return b
}

func (b *ShowMeasurementCardinality) OnDatabase(str string) *ShowMeasurementCardinality {
	b.q.Database = str
	return b
}

func (b *ShowMeasurementCardinality) From(sources ...*Measurement) *ShowMeasurementCardinality {
	for _, me := range sources {
		b.q.Sources = append(b.q.Sources, me.m)
	}
	return b
}

func (b *ShowMeasurementCardinality) GroupBy(fields ...GroupByIf) *ShowMeasurementCardinality {
	for _, f := range fields {
		if finalField := f.groupBy(); finalField != nil {
			b.q.Dimensions = append(b.q.Dimensions, finalField)
		}
	}
	return b
}

func (b *ShowMeasurementCardinality) Where(condition interface{}) *ShowMeasurementCardinality {
	switch condition := condition.(type) {
	case influxql.Expr:
		b.q.Condition = &influxql.ParenExpr{Expr: condition}
	case MathExprIf:
		b.q.Condition = condition.expr()
	}
	return b
}

func (b *ShowMeasurementCardinality) Build() (string, error) {
	return b.q.String(), nil
}

type ShowMeasurements struct {
	q *influxql.ShowMeasurementsStatement
}

func (b *ShowMeasurements) OnDatabase(db string) *ShowMeasurements {
	b.q.Database = db
	return b
}

func (b *ShowMeasurements) WithMeasurement(regexp *regexp.Regexp) *ShowMeasurements {
	b.q.Source = &influxql.Measurement{Regex: &influxql.RegexLiteral{Val: regexp}}
	return b
}

func (b *ShowMeasurements) Where(condition interface{}) *ShowMeasurements {
	switch condition := condition.(type) {
	case influxql.Expr:
		b.q.Condition = &influxql.ParenExpr{Expr: condition}
	case MathExprIf:
		b.q.Condition = condition.expr()
	}
	return b
}

func (b *ShowMeasurements) Limit(limit int) *ShowMeasurements {
	b.q.Limit = limit
	return b
}

func (b *ShowMeasurements) Offset(offset int) *ShowMeasurements {
	b.q.Offset = offset
	return b
}

func (b *ShowMeasurements) Build() (string, error) {
	return b.q.String(), nil
}

type ShowRetentionPolicies struct {
	db string
}

func (b *ShowRetentionPolicies) WithDatabse(db string) *ShowRetentionPolicies {
	b.db = db
	return b
}

func (b *ShowRetentionPolicies) Build() (string, error) {
	q := &influxql.ShowRetentionPoliciesStatement{Database: b.db}
	return q.String(), nil
}

type ShowSeries struct {
	q *influxql.ShowSeriesStatement
}

func (b *ShowSeries) OnDatabase(db string) *ShowSeries {
	b.q.Database = db
	return b
}

func (b *ShowSeries) WithMeasurement(sources ...*Measurement) *ShowSeries {
	for _, m := range sources {
		b.q.Sources = append(b.q.Sources, m.m)
	}
	return b
}

func (b *ShowSeries) Where(condition interface{}) *ShowSeries {
	switch condition := condition.(type) {
	case influxql.Expr:
		b.q.Condition = &influxql.ParenExpr{Expr: condition}
	case MathExprIf:
		b.q.Condition = condition.expr()
	}
	return b
}

func (b *ShowSeries) Limit(limit int) *ShowSeries {
	b.q.Limit = limit
	return b
}

func (b *ShowSeries) Offset(offset int) *ShowSeries {
	b.q.Offset = offset
	return b
}

func (b *ShowSeries) Build() (string, error) {
	return b.q.String(), nil
}

type ShowSeriesCadinality struct {
	q *influxql.ShowSeriesCardinalityStatement
}

func (b *ShowSeriesCadinality) Excat() *ShowSeriesCadinality {
	b.q.Exact = true
	return b
}

func (b *ShowSeriesCadinality) OnDatabase(db string) *ShowSeriesCadinality {
	b.q.Database = db
	return b
}

func (b *ShowSeriesCadinality) From(measuremnts ...*Measurement) *ShowSeriesCadinality {
	for _, m := range measuremnts {
		b.q.Sources = append(b.q.Sources, m.m)
	}
	return b
}

func (b *ShowSeriesCadinality) Where(condition interface{}) *ShowSeriesCadinality {
	switch condition := condition.(type) {
	case influxql.Expr:
		b.q.Condition = &influxql.ParenExpr{Expr: condition}
	case MathExprIf:
		b.q.Condition = condition.expr()
	}
	return b
}

func (b *ShowSeriesCadinality) Limit(limit int) *ShowSeriesCadinality {
	b.q.Limit = limit
	return b
}

func (b *ShowSeriesCadinality) Offset(offset int) *ShowSeriesCadinality {
	b.q.Offset = offset
	return b
}

func (b *ShowSeriesCadinality) GroupBy(fields ...GroupByIf) *ShowSeriesCadinality {
	for _, f := range fields {
		if finalField := f.groupBy(); finalField != nil {
			b.q.Dimensions = append(b.q.Dimensions, finalField)
		}
	}
	return b
}

func (b *ShowSeriesCadinality) Build() (string, error) {
	return b.q.String(), nil
}

type ShowStats struct {
	q *influxql.ShowStatsStatement
}

func (b *ShowStats) ForComponent(component string) *ShowStats {
	b.q.Module = component
	return b
}

func (b *ShowStats) ForIndexes() *ShowStats {
	return b.ForComponent("indexes")
}

func (b *ShowStats) Build() (string, error) {
	return b.q.String(), nil
}

type ShowTagKeyCardinality struct {
	q *influxql.ShowTagKeyCardinalityStatement
}

func (b *ShowTagKeyCardinality) Exact() *ShowTagKeyCardinality {
	b.q.Exact = true
	return b
}

func (b *ShowTagKeyCardinality) Limit(limit int) *ShowTagKeyCardinality {
	b.q.Limit = limit
	return b
}

func (b *ShowTagKeyCardinality) Offset(offset int) *ShowTagKeyCardinality {
	b.q.Offset = offset
	return b
}

func (b *ShowTagKeyCardinality) OnDatabase(str string) *ShowTagKeyCardinality {
	b.q.Database = str
	return b
}

func (b *ShowTagKeyCardinality) From(sources ...*Measurement) *ShowTagKeyCardinality {
	for _, me := range sources {
		b.q.Sources = append(b.q.Sources, me.m)
	}
	return b
}

func (b *ShowTagKeyCardinality) GroupBy(fields ...GroupByIf) *ShowTagKeyCardinality {
	for _, f := range fields {
		if finalField := f.groupBy(); finalField != nil {
			b.q.Dimensions = append(b.q.Dimensions, finalField)
		}
	}
	return b
}

func (b *ShowTagKeyCardinality) Where(condition interface{}) *ShowTagKeyCardinality {
	switch condition := condition.(type) {
	case influxql.Expr:
		b.q.Condition = &influxql.ParenExpr{Expr: condition}
	case MathExprIf:
		b.q.Condition = condition.expr()
	}
	return b
}

func (b *ShowTagKeyCardinality) Build() (string, error) {
	return b.q.String(), nil
}

type ShowTagKeys struct {
	q *influxql.ShowTagKeysStatement
}

func (b *ShowTagKeys) Limit(limit int) *ShowTagKeys {
	b.q.Limit = limit
	return b
}

func (b *ShowTagKeys) Offset(offset int) *ShowTagKeys {
	b.q.Offset = offset
	return b
}

func (b *ShowTagKeys) SLimit(limit int) *ShowTagKeys {
	b.q.SLimit = limit
	return b
}

func (b *ShowTagKeys) SOffset(offset int) *ShowTagKeys {
	b.q.SOffset = offset
	return b
}

func (b *ShowTagKeys) OnDatabase(str string) *ShowTagKeys {
	b.q.Database = str
	return b
}

func (b *ShowTagKeys) From(sources ...*Measurement) *ShowTagKeys {
	for _, me := range sources {
		b.q.Sources = append(b.q.Sources, me.m)
	}
	return b
}

func (b *ShowTagKeys) Where(condition interface{}) *ShowTagKeys {
	switch condition := condition.(type) {
	case influxql.Expr:
		b.q.Condition = &influxql.ParenExpr{Expr: condition}
	case MathExprIf:
		b.q.Condition = condition.expr()
	}
	return b
}

func (b *ShowTagKeys) Build() (string, error) {
	return b.q.String(), nil
}

type ShowTagValues struct {
	q *influxql.ShowTagValuesStatement
}

func (b *ShowTagValues) Limit(limit int) *ShowTagValues {
	b.q.Limit = limit
	return b
}

func (b *ShowTagValues) Offset(offset int) *ShowTagValues {
	b.q.Offset = offset
	return b
}

func (b *ShowTagValues) OnDatabase(str string) *ShowTagValues {
	b.q.Database = str
	return b
}

func (b *ShowTagValues) From(sources ...*Measurement) *ShowTagValues {
	for _, me := range sources {
		b.q.Sources = append(b.q.Sources, me.m)
	}
	return b
}

func (b *ShowTagValues) Where(condition interface{}) *ShowTagValues {
	switch condition := condition.(type) {
	case influxql.Expr:
		b.q.Condition = &influxql.ParenExpr{Expr: condition}
	case MathExprIf:
		b.q.Condition = condition.expr()
	}
	return b
}

// TODO : Rewrite with custom objects
func (b *ShowTagValues) WithTagKey(operator influxql.Token, tagKey influxql.Literal) *ShowTagValues {
	b.q.Op = operator
	b.q.TagKeyExpr = tagKey
	return b
}

func (b *ShowTagValues) Build() (string, error) {
	return b.q.String(), nil
}

type ShowTagValuesCardinality struct {
	q *influxql.ShowTagValuesCardinalityStatement
}

func (b *ShowTagValuesCardinality) WithTagKey(operator influxql.Token, tagKey influxql.Literal) *ShowTagValuesCardinality {
	b.q.Op = operator
	b.q.TagKeyExpr = tagKey
	return b
}

func (b *ShowTagValuesCardinality) Exact() *ShowTagValuesCardinality {
	b.q.Exact = true
	return b
}

func (b *ShowTagValuesCardinality) Limit(limit int) *ShowTagValuesCardinality {
	b.q.Limit = limit
	return b
}

func (b *ShowTagValuesCardinality) Offset(offset int) *ShowTagValuesCardinality {
	b.q.Offset = offset
	return b
}

func (b *ShowTagValuesCardinality) OnDatabase(str string) *ShowTagValuesCardinality {
	b.q.Database = str
	return b
}

func (b *ShowTagValuesCardinality) From(sources ...*Measurement) *ShowTagValuesCardinality {
	for _, me := range sources {
		b.q.Sources = append(b.q.Sources, me.m)
	}
	return b
}

func (b *ShowTagValuesCardinality) GroupBy(fields ...GroupByIf) *ShowTagValuesCardinality {
	for _, f := range fields {
		if finalField := f.groupBy(); finalField != nil {
			b.q.Dimensions = append(b.q.Dimensions, finalField)
		}
	}
	return b
}

func (b *ShowTagValuesCardinality) Where(condition interface{}) *ShowTagValuesCardinality {
	switch condition := condition.(type) {
	case influxql.Expr:
		b.q.Condition = &influxql.ParenExpr{Expr: condition}
	case MathExprIf:
		b.q.Condition = condition.expr()
	}
	return b
}

func (b *ShowTagValuesCardinality) Build() (string, error) {
	return b.q.String(), nil
}
