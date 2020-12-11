package influxqb

import (
	"fmt"
	"github.com/influxdata/influxql"
	"regexp"
	"strconv"
	"time"
)

type Order bool

const DESC Order = false
const ASC Order = true

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

type FieldIf interface {
	field() *influxql.Field
}
type Field struct {
	Name  string
	Alias string
	Type  influxql.DataType
}

func (f *Field) WithName(str string) *Field {
	f.Name = str
	return f
}

func (f *Field) WithAlias(str string) *Field {
	f.Alias = str
	return f
}

func (f *Field) WithType(inType influxql.DataType) *Field {
	f.Type = inType
	return f
}

func NewField(name string) *Field {
	return &Field{Name: name}
}

type Wildcard struct {
	FieldIf
}

func NewWildcardField() *Wildcard {
	return &Wildcard{}
}

type GroupByIf interface {
	groupBy() *influxql.Dimension
}

type MathExprIf interface {
	expr() *influxql.ParenExpr
}

type TimeSampling struct {
	Interval time.Duration
}

func NewTimeSampling(it time.Duration) *TimeSampling {
	return &TimeSampling{Interval: it}
}

type Function struct {
	Name  string
	Args  []interface{}
	Alias string
}

func (f *Function) WithAlias(alias string) *Function {
	f.Alias = alias
	return f
}

func (f *Function) WithArgs(args ...interface{}) *Function {
	f.Args = append(f.Args, args...)
	return f
}

func (f *Function) WithArg(arg interface{}) *Function {
	return f.WithArgs(arg)
}

func NewFunction(functionName string) *Function {
	return &Function{Name: functionName}
}

type MathExpr struct {
	Expr  string
	Alias string
}

func (m *MathExpr) WithAlias(alias string) *MathExpr {
	m.Alias = alias
	return m
}

func NewMathExpr(expr string) *MathExpr {
	return &MathExpr{Expr: expr}
}

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

func (w *Wildcard) field() *influxql.Field {
	return &influxql.Field{Expr: &influxql.Wildcard{}}
}

func (f *Field) field() *influxql.Field {
	return &influxql.Field{Expr: &influxql.VarRef{Val: f.Name, Type: f.Type}, Alias: f.Alias}
}

func (f *Field) groupBy() *influxql.Dimension {
	return &influxql.Dimension{Expr: f.field().Expr}
}

func (s *TimeSampling) groupBy() *influxql.Dimension {
	f := &Function{
		Name: "time",
		Args: []interface{}{s.Interval},
	}
	return &influxql.Dimension{Expr: f.field().Expr}
}

func (f *Function) field() *influxql.Field {
	args := []influxql.Expr{}

	for _, a := range f.Args {
		args = append(args, toExpr(a))
	}

	return &influxql.Field{
		Expr: &influxql.Call{
			Name: f.Name,
			Args: args,
		},
		Alias: f.Alias,
	}
}

func toExpr(a interface{}) influxql.Expr {
	switch a.(type) {
	case int:
		intval := int64(a.(int))
		return &influxql.IntegerLiteral{Val: intval}
	case int8:
		intval := int64(a.(int8))
		return &influxql.IntegerLiteral{Val: intval}
	case int16:
		intval := int64(a.(int16))
		return &influxql.IntegerLiteral{Val: intval}
	case int32:
		intval := int64(a.(int32))
		return &influxql.IntegerLiteral{Val: intval}
	case int64:
		intval := a.(int64)
		return &influxql.IntegerLiteral{Val: intval}
	case uint:
		uintval := int64(a.(uint))
		return &influxql.IntegerLiteral{Val: uintval}
	case uint8:
		intval := int64(a.(uint8))
		return &influxql.IntegerLiteral{Val: intval}
	case uint16:
		intval := int64(a.(uint16))
		return &influxql.IntegerLiteral{Val: intval}
	case uint32:
		intval := int64(a.(uint32))
		return &influxql.IntegerLiteral{Val: intval}
	case uint64:
		intval := int64(a.(uint64))
		return &influxql.IntegerLiteral{Val: intval}
	case float32:
		floatval := float64(a.(float32))
		return &influxql.NumberLiteral{Val: floatval}
	case float64:
		floatval := a.(float64)
		return &influxql.NumberLiteral{Val: floatval}
	case string:
		return &influxql.StringLiteral{Val: a.(string)}
	case time.Time:
		return &influxql.TimeLiteral{Val: a.(time.Time)}
	case *time.Time:
		tmp := *a.(*time.Time)
		return &influxql.TimeLiteral{Val: tmp}
	case time.Duration:
		return &influxql.DurationLiteral{Val: a.(time.Duration)}
	case *Field:
		f := a.(*Field)
		return &influxql.VarRef{Val: f.Name, Type: f.Type}
	case Field:
		field := a.(Field)
		return &influxql.VarRef{Val: field.Name, Type: field.Type}
	case *Parenthesis:
		t := a.(*Parenthesis).compute()
		return t
	case influxql.Expr:
		return a.(influxql.Expr)
	default:
		return nil
	}

}

type Parenthesis struct {
	Expr []interface{}
}

func (p *Parenthesis) findParenthesis() (int, int, *influxql.ParenExpr) {
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
					p2 := &Parenthesis{Expr: p.Expr[firtOpen+1 : i]}
					return firtOpen, i, p2.compute()
				}
			}
		case *Parenthesis:
			p2 := v.(*Parenthesis)
			return i, i, p2.compute()
		}

	}

	return -1, -1, nil

}

func (p *Parenthesis) findLevelOps(level int) (int, int, *influxql.BinaryExpr) {

	for i, v := range p.Expr {
		switch v.(type) {
		case influxql.Token:
			token := v.(influxql.Token)
			if token.Precedence() == level {

				LHS := toExpr(p.Expr[i-1])
				RHS := toExpr(p.Expr[i+1])

				return i - 1, i + 1, &influxql.BinaryExpr{
					Op:  token,
					LHS: LHS,
					RHS: RHS,
				}
			}
		}
	}

	return -1, -1, nil
}

func (p *Parenthesis) compute() *influxql.ParenExpr {

	//Remove trailing paren
	//if p.Expr[0] == influxql.LPAREN {
	//	p.Expr = p.Expr[1:]
	//}
	//
	//if p.Expr[len(p.Expr)-1] == influxql.RPAREN {
	//	p.Expr = p.Expr[:len(p.Expr)-1]
	//}

	//fmt.Println(p.Expr)

	//0. FIND ( )
	for {
		s, e, expr := p.findParenthesis()
		if s != -1 && e != -1 {
			if e+1 >= len(p.Expr) {
				//p2 := &Parenthesis{Expr: p.Expr[s+1 : e]}
				p.Expr = append(p.Expr[:s], expr)
			} else {
				//p2 := &Parenthesis{Expr: p.Expr[s+1 : e]}
				tmp := p.Expr[e+1:]
				p.Expr = append(p.Expr[:s], expr)
				p.Expr = append(p.Expr, tmp...)
			}
		} else {
			break
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
				//if e >= len(p.Expr) {
				//	//Should never go here it means that a token is at the end of the array
				//	p.Expr = append(p.Expr[:s], expr)
				//} else {
				if e < len(p.Expr) {
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

func NewMath() *Math {
	return &Math{Expr: []interface{}{}}
}

func (m *Math) WithAlias(alias string) *Math {
	m.Alias = alias
	return m
}

func (m *Math) WithExpr(elements ...interface{}) *Math {
	m.Expr = []interface{}{}
	m.Expr = append(m.Expr, elements...)
	return m
}

func (m *Math) field() *influxql.Field {

	v := &Parenthesis{m.Expr}
	return &influxql.Field{Expr: v.compute(), Alias: m.Alias}

}

func (m *Math) expr() *influxql.ParenExpr {
	p := Parenthesis{Expr: m.Expr}
	return p.compute()
}

func (m *MathExpr) toMath() *Math {
	expr, err := influxql.ParseExpr(m.Expr)

	if err != nil {
		fmt.Println("Something wrong in MathExpr.toMath()")
	}

	return &Math{Expr: []interface{}{expr}, Alias: m.Alias}
}

func (m *MathExpr) field() *influxql.Field {
	return &influxql.Field{Expr: m.toMath().expr(), Alias: m.Alias}
}

func (m *MathExpr) expr() *influxql.ParenExpr {
	return m.toMath().expr()
}

type SelectBuilder struct {
	selectStatement *influxql.SelectStatement
}

func (q *SelectBuilder) AddSelect(fields ...interface{}) *SelectBuilder {
	return q.Select(fields...)
}

func (q *SelectBuilder) Select(fields ...interface{}) *SelectBuilder {
	for _, f := range fields {
		switch f.(type) {
		case influxql.Expr:
			d := f.(influxql.Expr)
			q.selectStatement.Fields = append(q.selectStatement.Fields, &influxql.Field{Expr: &influxql.ParenExpr{Expr: d}})
			break
		case FieldIf:
			final := f.(FieldIf)
			if finalField := final.field(); finalField != nil {
				q.selectStatement.Fields = append(q.selectStatement.Fields, finalField)
			}
			break
		case *regexp.Regexp:
			final := f.(*regexp.Regexp)
			q.selectStatement.Fields = append(q.selectStatement.Fields, &influxql.Field{Expr: &influxql.RegexLiteral{Val: final}})
			break
		case string:
			if f.(string) == "*" {
				f2 := &Wildcard{}
				q.selectStatement.Fields = append(q.selectStatement.Fields, f2.field())
			} else {
				f2 := &Field{Name: f.(string)}
				q.selectStatement.Fields = append(q.selectStatement.Fields, f2.field())
			}

			break
		}
	}
	return q
}

func (q *SelectBuilder) AddGroupBy(fields ...GroupByIf) *SelectBuilder {
	return q.GroupBy(fields...)
}

func (q *SelectBuilder) GroupBy(fields ...GroupByIf) *SelectBuilder {
	for _, f := range fields {
		if finalField := f.groupBy(); finalField != nil {
			q.selectStatement.Dimensions = append(q.selectStatement.Dimensions, finalField)
		}
	}
	return q
}

func (q *SelectBuilder) AddFrom(fields ...string) *SelectBuilder {
	return q.From(fields...)
}

func (q *SelectBuilder) From(sources ...string) *SelectBuilder {
	for _, v := range sources {
		q.selectStatement.Sources = append(q.selectStatement.Sources, &influxql.Measurement{
			Name: v,
		})
	}
	return q
}

func (q *SelectBuilder) AddFromSubQuery(fields ...*SelectBuilder) *SelectBuilder {
	for _, qu := range fields {
		q.FromSubQuery(qu)
	}
	return q
}

func (q *SelectBuilder) FromSubQuery(subQuery *SelectBuilder) *SelectBuilder {
	q.selectStatement.Sources = append(q.selectStatement.Sources, &influxql.SubQuery{Statement: subQuery.selectStatement})
	return q
}

func (q *SelectBuilder) AddFromRegex(regexes ...*regexp.Regexp) *SelectBuilder {
	return q.FromRegex(regexes...)
}
func (q *SelectBuilder) FromRegex(regexes ...*regexp.Regexp) *SelectBuilder {
	for _, v := range regexes {
		q.selectStatement.Sources = append(q.selectStatement.Sources, &influxql.Measurement{Regex: &influxql.RegexLiteral{Val: v}})
	}
	return q
}

func (q *SelectBuilder) Fill(fillOption interface{}) *SelectBuilder {

	switch fillOption.(type) {
	case string:
		fillValue := FillFromStr(fillOption.(string))
		q.Fill(fillValue)
		break
	case FillOption:
		option, value := fillOption.(FillOption).get()
		q.selectStatement.Fill = option
		q.selectStatement.FillValue = value
		break
	case int:
		q.Fill(FillNumber{Value: float64(fillOption.(int))})
		break
	case float64:
		q.Fill(FillNumber{Value: fillOption.(float64)})
	default:
		q.Fill(FillNoFill{})
	}

	return q
}

func (q *SelectBuilder) Limit(limit int) *SelectBuilder {
	q.selectStatement.Limit = limit
	return q
}

func (q *SelectBuilder) SeriesLimit(limit int) *SelectBuilder {
	q.selectStatement.SLimit = limit
	return q
}

func (q *SelectBuilder) SeriesOffset(offset int) *SelectBuilder {
	q.selectStatement.SOffset = offset
	return q
}

func (q *SelectBuilder) Offset(offset int) *SelectBuilder {
	q.selectStatement.Offset = offset
	return q
}

func (q *SelectBuilder) Where(expr interface{}) *SelectBuilder {
	switch expr.(type) {
	case influxql.Expr:
		q.selectStatement.Condition = &influxql.ParenExpr{Expr: expr.(influxql.Expr)}
		break
	case MathExprIf:
		q.selectStatement.Condition = expr.(MathExprIf).expr()
		break
	}
	return q
}

func (q *SelectBuilder) Build() (string, error) {
	return q.selectStatement.String(), nil
}

func (q *SelectBuilder) WithTimeZone(timeZone interface{}) *SelectBuilder {

	var tz *time.Location
	switch timeZone.(type) {
	case *time.Location:
		tz = timeZone.(*time.Location)
		break
	case string:
		tz, _ = time.LoadLocation(timeZone.(string))
	}

	q.selectStatement.Location = tz
	return q
}

func (q *SelectBuilder) OrderBy(str string, order Order) *SelectBuilder {
	q.selectStatement.SortFields = append(q.selectStatement.SortFields, &influxql.SortField{Name: str, Ascending: bool(order)})
	return q
}

func (q *SelectBuilder) Into(mesurement *Measurement) *SelectBuilder {
	//TODO Improve this function
	q.selectStatement.Target = &influxql.Target{Measurement: mesurement.m}
	return q
}

func Or(LH interface{}, RH interface{}) *influxql.BinaryExpr {
	return &influxql.BinaryExpr{LHS: toExpr(LH), RHS: toExpr(RH), Op: influxql.OR}
}

func And(LH interface{}, RH interface{}) *influxql.BinaryExpr {
	return &influxql.BinaryExpr{LHS: toExpr(LH), RHS: toExpr(RH), Op: influxql.AND}
}

func Eq(LH interface{}, RH interface{}) *influxql.BinaryExpr {
	return &influxql.BinaryExpr{LHS: toExpr(LH), RHS: toExpr(RH), Op: influxql.EQ}
}

func LessThan(LH interface{}, RH interface{}) *influxql.BinaryExpr {
	return &influxql.BinaryExpr{LHS: toExpr(LH), RHS: toExpr(RH), Op: influxql.LT}
}

func GreaterThan(LH interface{}, RH interface{}) *influxql.BinaryExpr {
	return &influxql.BinaryExpr{LHS: toExpr(LH), RHS: toExpr(RH), Op: influxql.GT}
}

func LessThanEq(LH interface{}, RH interface{}) *influxql.BinaryExpr {
	return &influxql.BinaryExpr{LHS: toExpr(LH), RHS: toExpr(RH), Op: influxql.LTE}
}

func GreaterThanEq(LH interface{}, RH interface{}) *influxql.BinaryExpr {
	return &influxql.BinaryExpr{LHS: toExpr(LH), RHS: toExpr(RH), Op: influxql.GTE}
}

func Add(LH interface{}, RH interface{}) *influxql.BinaryExpr {
	return &influxql.BinaryExpr{LHS: toExpr(LH), RHS: toExpr(RH), Op: influxql.ADD}
}
func Subtract(LH interface{}, RH interface{}) *influxql.BinaryExpr {
	return &influxql.BinaryExpr{LHS: toExpr(LH), RHS: toExpr(RH), Op: influxql.SUB}
}

func Divide(LH interface{}, RH interface{}) *influxql.BinaryExpr {
	return &influxql.BinaryExpr{LHS: toExpr(LH), RHS: toExpr(RH), Op: influxql.DIV}
}

func Multiply(LH interface{}, RH interface{}) *influxql.BinaryExpr {
	return &influxql.BinaryExpr{LHS: toExpr(LH), RHS: toExpr(RH), Op: influxql.MUL}
}

func Modulus(LH interface{}, RH interface{}) *influxql.BinaryExpr {
	return &influxql.BinaryExpr{LHS: toExpr(LH), RHS: toExpr(RH), Op: influxql.MOD}
}

func NotEq(LH interface{}, RH interface{}) *influxql.BinaryExpr {
	return &influxql.BinaryExpr{LHS: toExpr(LH), RHS: toExpr(RH), Op: influxql.NEQ}
}

func FillFromStr(str string) FillOption {

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
