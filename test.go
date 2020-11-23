package main

import (
	"fmt"
	"github.com/influxdata/influxql"
	"time"
)

func main() {

	//s := influxql.StringLiteral{Val: "TestField"}

	InfluxQlWrapper()
	//influxQLTest()
}

func InfluxQlWrapper() {

	dur, _ := time.ParseDuration("15h")

	builder := NewQueryBuilder()
	builder.Select(
		&Function{Name: "MEAN", Args: []interface{}{"colomn", time.Now(), 45.36, dur}},
		&Field{Name: "MyField"},
		&Math{Expr: []interface{}{
			"COL", influxql.ADD, 156, influxql.MUL, influxql.LPAREN, "COL2", influxql.MOD, 12, influxql.RPAREN,
		}},
	)
	builder.From("XTC_OLD'sk")
	//builder.FromRegex(regexp.MustCompile(`^\d+(ns|u|ms|s|m|h|d|w)?$`))
	builder.GroupBy(
		&Field{Name: "GroupByField"},
		&TimeSampling{interval: dur},
	)
	builder.Fill(FillNumber{45})
	builder.Limit(250)
	builder.Offset(15)
	builder.SeriesLimit(2)
	builder.SeriesOffset(8)

	//TODO: Add

	fmt.Println(builder.Build())
}

func influxQLTest() {

	fiels := influxql.Fields{}

	fiels = append(fiels, &influxql.Field{Expr: &influxql.StringLiteral{Val: "Str'ingVa\"lue"}, Alias: "COOL"})
	fiels = append(fiels, &influxql.Field{Expr: &influxql.Wildcard{}})

	fiels = append(fiels, &influxql.Field{Expr: &influxql.Call{Name: "mean", Args: []influxql.Expr{&influxql.StringLiteral{Val: "ej';dij'dfie\""}}}, Alias: "WILM"})

	fiels = append(fiels, &influxql.Field{Expr: &influxql.BinaryExpr{
		Op:  influxql.ADD,
		LHS: &influxql.NumberLiteral{Val: 351},
		RHS: &influxql.StringLiteral{Val: "CA"},
	}})

	fiels = append(fiels, &influxql.Field{Expr: &influxql.TimeLiteral{Val: time.Now()}})

	meas := &influxql.Measurement{
		//Database: "DB",
		Name: "MEAS",
	}

	sources := influxql.Sources{}
	sources = append(sources, meas)

	q := influxql.SelectStatement{
		Fields:     fiels,
		Target:     nil,
		Dimensions: nil,
		Sources:    sources,
		Condition:  nil,
		SortFields: nil,
		Limit:      0,
		Offset:     0,
		SLimit:     0,
		SOffset:    0,
		IsRawQuery: false,
		Fill:       0,
		FillValue:  nil,
		Location:   nil,
		TimeAlias:  "",
		OmitTime:   false,
		StripName:  false,
		EmitName:   "",
		Dedupe:     false,
	}

	fmt.Println(q.String())

	fmt.Println("Hello")
}
