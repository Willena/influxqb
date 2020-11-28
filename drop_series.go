package influxqb

import (
	"github.com/influxdata/influxql"
)

type DropSeriesBuilder struct {
	del *influxql.DropSeriesStatement
}

func (b *DropSeriesBuilder) From(sources ...*Measurement) *DropSeriesBuilder {
	for _, f := range sources {
		b.del.Sources = append(b.del.Sources, f.m)
	}
	return b
}

func (b *DropSeriesBuilder) Where(condition interface{}) *DropSeriesBuilder {
	switch condition.(type) {
	case influxql.Expr:
		b.del.Condition = &influxql.ParenExpr{Expr: condition.(influxql.Expr)}
		break
	case MathExprIf:
		b.del.Condition = condition.(MathExprIf).expr()
		break
	}
	return b
}

func (b *DropSeriesBuilder) Build() (string, error) {
	return b.del.String(), nil
}
