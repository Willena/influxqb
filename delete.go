package influxqb

import (
	"github.com/influxdata/influxql"
	"regexp"
)

type DeleteBuilder struct {
	del *influxql.DeleteStatement
}

func (b *DeleteBuilder) From(source string) *DeleteBuilder {
	b.del.Source = &influxql.Measurement{Name: source}
	return b
}

func (b *DeleteBuilder) Where(condition interface{}) *DeleteBuilder {
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

func (b *DeleteBuilder) Build() (string, error) {

	if b.del.Source == nil {
		b.del.Source = &influxql.Measurement{
			Regex: &influxql.RegexLiteral{Val: regexp.MustCompile(".*")},
		}
	}

	return b.del.String(), nil
}
