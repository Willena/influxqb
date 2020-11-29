package influxqb

import (
	"github.com/influxdata/influxql"
)

type ExplainBuilder struct {
	explainStatement *influxql.ExplainStatement
	selectBuilder    *SelectBuilder
}

func (b *ExplainBuilder) Analyze() *ExplainBuilder {
	b.explainStatement.Analyze = true
	return b
}

func (b *ExplainBuilder) WithSelectBuilder(builder *SelectBuilder) *ExplainBuilder {
	b.selectBuilder = builder
	return b
}

func (b *ExplainBuilder) SelectBuilder() *SelectBuilder {
	return b.selectBuilder
}

func (b *ExplainBuilder) Build() (string, error) {
	b.explainStatement.Statement = b.selectBuilder.selectStatement
	return b.explainStatement.String(), nil
}
