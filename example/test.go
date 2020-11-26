package main

import (
	"fmt"
	"github.com/influxdata/influxql"
	"github.com/willena/influxql-query-builder"
	"time"
)

func main() {

	//s := influxql.StringLiteral{Val: "TestField"}

	//s, err := influxql.ParseStatement("SELECT toto, titi FROM a WHERE a = 'r'; DROP MEASUREMENT kkk'")
	//
	//
	//fmt.Println(s)
	//fmt.Println(err)
	//
	//s2 := influxql.MustParseExpr("A >= rt',;ja AND T < 15h")
	//fmt.Println(s2)

	q, err := influxql.ParseStatement("SELECT \"toto\" FROM \"uuu\"")

	f := influxql.Field{Expr: &influxql.VarRef{Val: "Tot\"o", Type: influxql.String}}

	fmt.Println(f.String())
	fmt.Println(q, err)
	fmt.Println(&influxql.Measurement{Name: "r\"fjf"})

	//InfluxQlWrapper()
	//influxQLTest()
}

func InfluxQlWrapper() {

	dur, _ := time.ParseDuration("15h")

	builder := influxqb.NewQueryBuilder()
	builder.Select(
		&influxqb.Function{Name: "MEAN", Args: []interface{}{"colomn", time.Now(), 45.36, dur}},
		&influxqb.Field{Name: "MyField"},
	)
	builder.From("XTC_OLD'sk")
	//builder.FromRegex(regexp.MustCompile(`^\d+(ns|u|ms|s|m|h|d|w)?$`))
	builder.GroupBy(
		&influxqb.Field{Name: "GroupByField"},
		&influxqb.TimeSampling{Interval: dur},
	)
	builder.Fill(influxqb.FillNumber{45})
	builder.Limit(250)
	builder.Offset(15)
	builder.SeriesLimit(2)
	builder.SeriesOffset(8)
	builder.Where(
		&influxqb.Math{Expr: []interface{}{influxqb.Field{Name: "Tptp"}, influxql.EQ, "data", influxql.AND, "ooo", influxql.EQ, 16.55}})

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
