package main

import (
	"fmt"
	"github.com/influxdata/influxql"
	"regexp"
	"time"
)

type FillOption interface {
	get() (influxql.FillOption, interface{})
}

type FillNoFill struct {}
type FillNumber struct {
	value int
}
type FillNull struct {}
type FillPrevious struct {}
type FillLinear struct {}

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
	return influxql.NumberFill, receiver.value
}






type FieldIf interface {
	field() *influxql.Field
}

type GroupByIf interface {
	groupBy() *influxql.Dimension
}

type Field struct {
	FieldIf
	Name  string
	Alias string
}

func (f *Field) field() *influxql.Field {
	return &influxql.Field{Expr: &influxql.StringLiteral{Val: f.Name}, Alias: f.Alias}
}

func (f *Field) groupBy() *influxql.Dimension {
	return &influxql.Dimension{Expr: f.field().Expr}
}

type TimeSampling struct {
	GroupByIf
	interval time.Duration
}

func (s *TimeSampling) groupBy() *influxql.Dimension {
	f := &Function{
		Name: "time",
		Args: []interface{}{s.interval},
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

type Math struct {
	Expr  []interface{}
	Alias string
}

func (m *Math) field() *influxql.Field {

	//TODO implement the logic to read the array
	//     for now it is only reading the first native binary Expr

	for _, v := range m.Expr {
		switch v.(type) {
		case *influxql.BinaryExpr:
			return &influxql.Field{Expr: v.(*influxql.BinaryExpr), Alias: m.Alias}
		}
	}

	fmt.Println("Warning: Math type only accept a single *influxql.BinaryExpr for the moment.")

	return nil
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

func (q *QueryBuilder) Build() string {
	return q.selectStatement.String()
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{selectStatement: influxql.SelectStatement{
		Fields:  []*influxql.Field{},
		Sources: []influxql.Source{},
	}}
}
