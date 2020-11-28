package influxqb

import "github.com/influxdata/influxql"

type DropRetentionPolicyBuilder struct {
	dcq *influxql.DropRetentionPolicyStatement
}

func (b *DropRetentionPolicyBuilder) WithRetentionPolicy(RetentionPolicy string) *DropRetentionPolicyBuilder {
	b.dcq.Name = RetentionPolicy
	return b
}

func (b *DropRetentionPolicyBuilder) WithDatabase(database string) *DropRetentionPolicyBuilder {
	b.dcq.Database = database
	return b
}

func (b *DropRetentionPolicyBuilder) Build() (string, error) {
	return b.dcq.String(), nil
}
