package influxqb

import (
	"github.com/influxdata/influxql"
)

type DropSeriesBuilder struct {
	dss *influxql.DropSeriesStatement
}

func (b *DropSeriesBuilder) From(sources ...*Measurement) *DropSeriesBuilder {
	for _, f := range sources {
		b.dss.Sources = append(b.dss.Sources, f.m)
	}
	return b
}

func (b *DropSeriesBuilder) Where(condition interface{}) *DropSeriesBuilder {
	switch condition.(type) {
	case influxql.Expr:
		b.dss.Condition = &influxql.ParenExpr{Expr: condition.(influxql.Expr)}
		break
	case MathExprIf:
		b.dss.Condition = condition.(MathExprIf).expr()
		break
	}
	return b
}

func (b *DropSeriesBuilder) Build() (string, error) {
	return b.dss.String(), nil
}
