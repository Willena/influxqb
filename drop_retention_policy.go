package influxqb

import "github.com/influxdata/influxql"

type DropRetentionPolicyBuilder struct {
	drp *influxql.DropRetentionPolicyStatement
}

func (b *DropRetentionPolicyBuilder) WithRetentionPolicy(RetentionPolicy string) *DropRetentionPolicyBuilder {
	b.drp.Name = RetentionPolicy
	return b
}

func (b *DropRetentionPolicyBuilder) WithDatabase(database string) *DropRetentionPolicyBuilder {
	b.drp.Database = database
	return b
}

func (b *DropRetentionPolicyBuilder) Build() (string, error) {
	return b.drp.String(), nil
}
