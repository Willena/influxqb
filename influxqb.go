package influxqb

import "github.com/influxdata/influxql"

type BuilderIf interface {
	Build() (string, error)
}

func NewSelectBuilder() *SelectBuilder {
	return &SelectBuilder{selectStatement: &influxql.SelectStatement{}}
}

func NewAlterRetentionPolicyBuilder() *AlterRetentionPolicyBuilder {
	return &AlterRetentionPolicyBuilder{alterStm: &influxql.AlterRetentionPolicyStatement{}}
}