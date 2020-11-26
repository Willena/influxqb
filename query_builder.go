package influxqb

import (
	"fmt"
	"github.com/influxdata/influxql"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ()

type FillOption interface {
	get() (influxql.FillOption, interface{})
}

type FillNoFill struct{}
type FillNumber struct {
	Value float64
}
type FillNull struct{}
type FillPrevious struct{}
type FillLinear struct{}

func (receiver FillNoFill) get() (influxql.FillOption, interface{}) {
	return influxql.NoFill, nil
}
func (receiver FillNull) get() (influxql.FillOption, interface{}) {
	return influxql.NullFill, nil
}
func (receiver FillLinear) get() (influxql.FillOption, interface{}) {
	return influxql.LinearFill, nil
}
func (receiver FillPrevious) get() (influxql.FillOption, interface{}) {
	return influxql.PreviousFill, nil
}
func (receiver FillNumber) get() (influxql.FillOption, interface{}) {
	return influxql.NumberFill, receiver.Value
}

type FieldIf interface {
	field() *influxql.Field
}

type GroupByIf interface {
	groupBy() *influxql.Dimension
}

type MathExprIf interface {
	expr() *influxql.ParenExpr
}

type Wildcard struct {
	FieldIf
}

func (w *Wildcard) field() *influxql.Field {
	return &influxql.Field{Expr: &influxql.Wildcard{}}
}

type Field struct {
	FieldIf
	Name  string
	Alias string
	Type  influxql.DataType
}

func (f *Field) field() *influxql.Field {
	return &influxql.Field{Expr: &influxql.VarRef{Val: f.Name, Type: f.Type}, Alias: f.Alias}
}

func (f *Field) groupBy() *influxql.Dimension {
	return &influxql.Dimension{Expr: f.field().Expr}
}

type TimeSampling struct {
	GroupByIf
	Interval time.Duration
}

func (s *TimeSampling) groupBy() *influxql.Dimension {
	f := &Function{
		Name: "time",
		Args: []interface{}{s.Interval},
	}
	return &influxql.Dimension{Expr: f.field().Expr}
}

type Function struct {
	FieldIf
	Name  string
	Args  []interface{}
	Alias string
}

func (f *Function) field() *influxql.Field {
	args := []influxql.Expr{}

	for _, a := range f.Args {
		var floatval float64
		switch a.(type) {
		case int:
			floatval = float64(a.(int))
		case int8:
			floatval = float64(a.(int8))
		case int16:
			floatval = float64(a.(int16))
		case int32:
			floatval = float64(a.(int32))
		case int64:
			floatval = float64(a.(int64))
		case uint:
			floatval = float64(a.(uint))
		case uint8:
			floatval = float64(a.(uint8))
		case uint16:
			floatval = float64(a.(uint16))
		case uint32:
			floatval = float64(a.(uint32))
		case uint64:
			floatval = float64(a.(uint64))
		case float32:
			floatval = float64(a.(float32))
		case float64:
			floatval = a.(float64)
			args = append(args, &influxql.NumberLiteral{Val: floatval})
			break
		case string:
			args = append(args, &influxql.StringLiteral{Val: a.(string)})
		case time.Time:
			args = append(args, &influxql.TimeLiteral{Val: a.(time.Time)})
		case time.Duration:
			args = append(args, &influxql.DurationLiteral{Val: a.(time.Duration)})
		case influxql.Expr:
			args = append(args, a.(influxql.Expr))
		case *Field:
			f := a.(*Field)
			args = append(args, &influxql.VarRef{Val: f.Name, Type: f.Type})
		default:
			//If not a number or string or time
		}
	}

	return &influxql.Field{
		Expr: &influxql.Call{
			Name: f.Name,
			Args: args,
		},
		Alias: f.Alias,
	}
}

func toExpr(a interface{}) []interface{} {
	ret := make([]interface{}, 0)

	var floatval float64
	switch a.(type) {
	case int:
		floatval = float64(a.(int))
		return append(ret, &influxql.NumberLiteral{Val: floatval})
	case int8:
		floatval = float64(a.(int8))
		return append(ret, &influxql.NumberLiteral{Val: floatval})
	case int16:
		floatval = float64(a.(int16))
		return append(ret, &influxql.NumberLiteral{Val: floatval})
	case int32:
		floatval = float64(a.(int32))
		return append(ret, &influxql.NumberLiteral{Val: floatval})
	case int64:
		floatval = float64(a.(int64))
		return append(ret, &influxql.NumberLiteral{Val: floatval})
	case uint:
		floatval = float64(a.(uint))
		return append(ret, &influxql.NumberLiteral{Val: floatval})
	case uint8:
		floatval = float64(a.(uint8))
		return append(ret, &influxql.NumberLiteral{Val: floatval})
	case uint16:
		floatval = float64(a.(uint16))
		return append(ret, &influxql.NumberLiteral{Val: floatval})
	case uint32:
		floatval = float64(a.(uint32))
		return append(ret, &influxql.NumberLiteral{Val: floatval})
	case uint64:
		floatval = float64(a.(uint64))
		return append(ret, &influxql.NumberLiteral{Val: floatval})
	case float32:
		floatval = float64(a.(float32))
		return append(ret, &influxql.NumberLiteral{Val: floatval})
	case float64:
		floatval = a.(float64)
		return append(ret, &influxql.NumberLiteral{Val: floatval})
	case string:
		return append(ret, &influxql.StringLiteral{Val: a.(string)})
	case time.Time:
		return append(ret, &influxql.TimeLiteral{Val: a.(time.Time)})
	case *time.Time:
		var tmp time.Time
		tmp = (*(a.(*time.Time)))
		return append(ret, &influxql.TimeLiteral{Val: tmp})
	case time.Duration:
		return append(ret, &influxql.DurationLiteral{Val: a.(time.Duration)})
	case Field:
		field := a.(Field)
		return append(ret, &influxql.VarRef{Val: field.Name, Type: field.Type})
	case Parenthesis:
		t := a.(*Parenthesis).compute()
		return append(ret, &influxql.ParenExpr{Expr: t})
	case influxql.Expr:
		return append(ret, a)
	default:
		return ret
	}

}

type Parenthesis struct {
	Expr []interface{}
}

func (p *Parenthesis) findParenthesis() (int, int) {
	parenOpened := 0
	firtOpen := -1
	//lastClose := - 1
	for i, v := range p.Expr {
		//fmt.Print(v)
		switch v.(type) {
		case influxql.Token:
			token := v.(influxql.Token)
			if token == influxql.LPAREN {
				if firtOpen == -1 {
					firtOpen = i
				}
				parenOpened++
			} else if token == influxql.RPAREN {
				parenOpened--
				if parenOpened < 0 {
					fmt.Println("Math expression not valid")
				}

				if parenOpened == 0 {
					return firtOpen, i
				}

			}
		case *Parenthesis:
			v.(*Parenthesis).compute()
		}

	}

	return -1, -1

}

func (p *Parenthesis) findLevelOps(level int) (int, int, *influxql.BinaryExpr) {

	for i, v := range p.Expr {
		switch v.(type) {
		case influxql.Token:
			token := v.(influxql.Token)
			if token.Precedence() == level {
				//fmt.Print("Token : ", v, " Level", level)
				//fmt.Println(toExpr(p.Expr[i-1]))
				LHS := toExpr(p.Expr[i-1])
				RHS := toExpr(p.Expr[i+1])

				//fmt.Println(LHS, RHS)

				return i - 1, i + 1, &influxql.BinaryExpr{
					Op:  token,
					LHS: LHS[0].(influxql.Expr),
					RHS: RHS[0].(influxql.Expr),
				}
			}
		case *Parenthesis:
			v.(*Parenthesis).compute()
		}
	}

	return -1, -1, nil
}

func (p *Parenthesis) compute() *influxql.ParenExpr {

	//Remove trailing paren
	if p.Expr[0] == influxql.LPAREN {
		p.Expr = p.Expr[1:]
	}

	if p.Expr[len(p.Expr)-1] == influxql.RPAREN {
		p.Expr = p.Expr[:len(p.Expr)-1]
	}

	//0. FIND ( )
	s, e := p.findParenthesis()
	if s != -1 && e != -1 {
		if e >= len(p.Expr) {
			p2 := &Parenthesis{Expr: p.Expr[s:e]}
			p.Expr = append(p.Expr[:s], p2.compute())
		} else {
			p2 := &Parenthesis{Expr: p.Expr[s : e+1]}
			tmp := p.Expr[e+1:]
			p.Expr = append(p.Expr[:s], p2.compute())
			p.Expr = append(p.Expr, tmp...)
		}

	}

	//5. FIND  *  /  %  <<  >>  &  &^
	//4. FIND   +  -  |  ^
	//3. FIND  ==  !=  <  <=  >  >=
	//2. FIND &&
	//1. FIND ||

	for i := 5; i >= 0; i-- {
		for {
			s, e, expr := p.findLevelOps(i)
			if s != -1 && e != -1 {
				if e >= len(p.Expr) {
					p.Expr = append(p.Expr[:s], expr)
				} else {
					tmp := p.Expr[e+1:]
					p.Expr = append(p.Expr[:s], expr)
					p.Expr = append(p.Expr, tmp...)
				}
			} else {
				break
			}
		}
	}

	//fmt.Println(p.Expr)
	//6. Returns
	//Should be single element
	return &influxql.ParenExpr{Expr: p.Expr[0].(influxql.Expr)}
}

type Math struct {
	Expr  []interface{}
	Alias string
}

func (m *Math) field() *influxql.Field {

	v := &Parenthesis{m.Expr}
	return &influxql.Field{Expr: v.compute(), Alias: m.Alias}

}

func (m *Math) expr() *influxql.ParenExpr {
	p := Parenthesis{Expr: m.Expr}
	//fmt.Println("EXPR", m.Expr)
	return p.compute()
}

type MathExpr struct {
	Expr  string
	Alias string
}

func (m *MathExpr) tokenFromStr(str string) influxql.Token {
	switch str {
	case ">=":
		return influxql.GTE
	case ">":
		return influxql.GT
	case "<=":
		return influxql.LTE
	case "<":
		return influxql.LT
	case "&&":
		return influxql.AND
	case "&":
		return influxql.BITWISE_AND
	case "||":
		return influxql.OR
	case "|":
		return influxql.BITWISE_OR
	case "=":
		return influxql.EQ
	case "!=":
		return influxql.NEQ
	case "*":
		return influxql.MUL
	case "/":
		return influxql.DIV
	case "%":
		return influxql.MOD
	case "+":
		return influxql.ADD
	case "-":
		return influxql.SUB
	case "^":
		return influxql.BITWISE_XOR
	case "(":
		return influxql.LPAREN
	case ")":
		return influxql.RPAREN
	default:
		return -1
	}
}

func ParseMathExpr(str string) (*influxql.Expr, error) {
	expr, err := influxql.ParseExpr(str)
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (m *MathExpr) toMath() *Math {
	var re = regexp.MustCompile(`(?m)>=|>|<=|<|&&|&|\|\||\||!=|=|\*|\/|%|\+|-|\^|\(|\)`)
	var array []interface{}
	m.Expr = strings.Replace(m.Expr, " ", "", -1)
	lastI := 0

	for _, matches := range re.FindAllStringIndex(m.Expr, -1) {
		token := m.Expr[matches[0]:matches[1]]

		if matches[0]-lastI == 0 && lastI != 0 {
		} else {
			array = append(array, m.Expr[lastI:matches[0]])
		}

		array = append(array, m.tokenFromStr(token))
		lastI = matches[1]
	}

	array = append(array, m.Expr[lastI:])

	for i, v := range array {
		switch v.(type) {
		case influxql.Token:
			continue
		case string:
			//Read int
			v2, err := strconv.ParseFloat(v.(string), 64)
			if err == nil {
				array[i] = v2
				continue
			}

			timeDur, err := time.ParseDuration(v.(string))
			if err == nil {
				array[i] = timeDur
			}
		}
	}

	return &Math{Expr: array, Alias: m.Alias}
}

func (m *MathExpr) field() *influxql.Field {
	return &influxql.Field{Expr: m.toMath().expr(), Alias: m.Alias}
}

func (m *MathExpr) expr() *influxql.ParenExpr {
	return m.toMath().expr()
}

type QueryBuilder struct {
	selectStatement influxql.SelectStatement
}

func (q *QueryBuilder) Select(fields ...FieldIf) *QueryBuilder {
	for _, f := range fields {
		if finalField := f.field(); finalField != nil {
			q.selectStatement.Fields = append(q.selectStatement.Fields, finalField)
		}
	}
	return q
}

func (q *QueryBuilder) GroupBy(fields ...GroupByIf) *QueryBuilder {
	for _, f := range fields {
		if finalField := f.groupBy(); finalField != nil {
			q.selectStatement.Dimensions = append(q.selectStatement.Dimensions, finalField)
		}
	}
	return q
}

func (q *QueryBuilder) From(sources ...string) *QueryBuilder {
	for _, v := range sources {
		q.selectStatement.Sources = append(q.selectStatement.Sources, &influxql.Measurement{
			Name: v,
		})
	}
	return q
}

func (q *QueryBuilder) Fill(fillOption FillOption) *QueryBuilder {
	option, value := fillOption.get()
	q.selectStatement.Fill = option
	q.selectStatement.FillValue = value
	return q
}

func (q *QueryBuilder) Limit(limit int) *QueryBuilder {
	q.selectStatement.Limit = limit
	return q
}

func (q *QueryBuilder) SeriesLimit(limit int) *QueryBuilder {
	q.selectStatement.SLimit = limit
	return q
}

func (q *QueryBuilder) SeriesOffset(offset int) *QueryBuilder {
	q.selectStatement.SOffset = offset
	return q
}

func (q *QueryBuilder) Offset(offset int) *QueryBuilder {
	q.selectStatement.Offset = offset
	return q
}

func (q *QueryBuilder) FromRegex(regexes ...*regexp.Regexp) *QueryBuilder {
	for _, v := range regexes {
		q.selectStatement.Sources = append(q.selectStatement.Sources, &influxql.Measurement{Regex: &influxql.RegexLiteral{Val: v}})
	}
	return q
}

func (q *QueryBuilder) Where(mathExpr MathExprIf) *QueryBuilder {
	q.selectStatement.Condition = mathExpr.expr()
	return q
}

func (q *QueryBuilder) Build() string {
	return q.selectStatement.String()
}

func FillOptionFromStr(str string) FillOption {

	switch str {
	case "none":
		return FillNoFill{}
	case "null":
		return FillNull{}
	case "linear":
		return FillLinear{}
	case "previous":
		return FillPrevious{}
	default:
		value, err := strconv.ParseFloat(str, 64)
		if err == nil {
			return FillNumber{Value: value}
		} else {
			return FillNoFill{}
		}
	}
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{selectStatement: influxql.SelectStatement{
		Fields:  []*influxql.Field{},
		Sources: []influxql.Source{},
	}}
}
