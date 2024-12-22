package influxqb

import (
	"regexp"

	"github.com/influxdata/influxql"
)

type DeleteBuilder struct {
	del *influxql.DeleteStatement
}

func (b *DeleteBuilder) From(source string) *DeleteBuilder {
	b.del.Source = &influxql.Measurement{Name: source}
	return b
}

func (b *DeleteBuilder) Where(condition interface{}) *DeleteBuilder {
	switch condition := condition.(type) {
	case influxql.Expr:
		b.del.Condition = &influxql.ParenExpr{Expr: condition}
	case MathExprIf:
		b.del.Condition = condition.expr()
	}
	return b
}

func (b *DeleteBuilder) Build() (string, error) {

	if b.del.Source == nil {
		b.del.Source = &influxql.Measurement{
			Regex: &influxql.RegexLiteral{Val: regexp.MustCompile(".*")},
		}
	}

	return b.del.String(), nil
}
