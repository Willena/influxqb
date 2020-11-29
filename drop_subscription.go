package influxqb

import (
	"github.com/influxdata/influxql"
)

type DropSubscriptionBuilder struct {
	dss *influxql.DropSubscriptionStatement
}

func (b *DropSubscriptionBuilder) WithName(name string) *DropSubscriptionBuilder {
	b.dss.Name = name
	return b
}

func (b *DropSubscriptionBuilder) WithDatabase(db string) *DropSubscriptionBuilder {
	b.dss.Database = db
	return b
}

func (b *DropSubscriptionBuilder) WithRetentionPolicy(policy string) *DropSubscriptionBuilder {
	b.dss.RetentionPolicy = policy
	return b
}

func (b *DropSubscriptionBuilder) Build() (string, error) {
	return b.dss.String(), nil
}
