package influxqb

import "github.com/influxdata/influxql"

type CreateSubscriptionBuilder struct {
	subStm *influxql.CreateSubscriptionStatement
}

type SubscriptionMode string

const ANY SubscriptionMode = "ANY"
const ALL SubscriptionMode = "ALL"

func (b *CreateSubscriptionBuilder) WithName(str string) *CreateSubscriptionBuilder {
	b.subStm.Name = str
	return b
}

func (b *CreateSubscriptionBuilder) WithDatabase(str string) *CreateSubscriptionBuilder {
	b.subStm.Database = str
	return b
}

func (b *CreateSubscriptionBuilder) WithRetentionPolicy(str string) *CreateSubscriptionBuilder {
	b.subStm.RetentionPolicy = str
	return b
}

func (b *CreateSubscriptionBuilder) WithMode(mode SubscriptionMode) *CreateSubscriptionBuilder {
	b.subStm.Mode = string(mode)
	return b
}

func (b *CreateSubscriptionBuilder) WithDestination(str string) *CreateSubscriptionBuilder {
	b.subStm.Destinations = append(b.subStm.Destinations, str)
	return b
}

func (b *CreateSubscriptionBuilder) WithDestinations(str ...string) *CreateSubscriptionBuilder {
	b.subStm.Destinations = append(b.subStm.Destinations, str...)
	return b
}

func (b *CreateSubscriptionBuilder) Build() (string, error) {
	return b.subStm.String(), nil
}
