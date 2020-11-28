package influxqb

import (
	"github.com/influxdata/influxql"
	"time"
)

type CreateContinuousQueryBuilder struct {
	continuousQuery *influxql.CreateContinuousQueryStatement
	selectBuilder   *SelectBuilder
}

func (b *CreateContinuousQueryBuilder) WithName(str string) *CreateContinuousQueryBuilder {
	b.continuousQuery.Name = str
	return b
}

func (b *CreateContinuousQueryBuilder) WithDatabase(str string) *CreateContinuousQueryBuilder {
	b.continuousQuery.Database = str
	return b
}

func (b *CreateContinuousQueryBuilder) WithSelectBuilder(builder *SelectBuilder) *CreateContinuousQueryBuilder {
	b.selectBuilder = builder
	return b
}

func (b *CreateContinuousQueryBuilder) WithResamplingInterval(interval time.Duration) *CreateContinuousQueryBuilder {
	b.continuousQuery.ResampleEvery = interval
	return b
}

func (b *CreateContinuousQueryBuilder) WithResamplingIntervalFromString(str string) *CreateContinuousQueryBuilder {
	interval, _ := time.ParseDuration(str)
	b.WithResamplingInterval(interval)
	return b
}

func (b *CreateContinuousQueryBuilder) WithTimeout(interval time.Duration) *CreateContinuousQueryBuilder {
	b.continuousQuery.ResampleFor = interval
	return b
}

func (b *CreateContinuousQueryBuilder) WithTimeoutString(str string) *CreateContinuousQueryBuilder {
	interval, _ := time.ParseDuration(str)
	b.WithTimeout(interval)
	return b
}
func (b *CreateContinuousQueryBuilder) SelectBuilder() *SelectBuilder {
	return b.selectBuilder
}

func (b *CreateContinuousQueryBuilder) Build() (string, error) {
	b.continuousQuery.Source = b.selectBuilder.selectStatement
	return b.continuousQuery.String(), nil
}
