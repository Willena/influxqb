package test

import (
	"fmt"
	"github.com/influxdata/influxql"
	"github.com/stretchr/testify/assert"
	"github.com/willena/influxqb"
	"regexp"
	"testing"
	"time"
)

var tz, err = time.LoadLocation("Europe/Paris")

var testSamples = []struct {
	d string
	b influxqb.BuilderIf
	s string
	e bool
}{
	{
		"Simple select statement",
		influxqb.NewSelectBuilder().Select("FieldA").From("MyMeasurement"),
		"SELECT FieldA FROM MyMeasurement",
		false,
	},
	{
		"Simple select * statement",
		influxqb.NewSelectBuilder().Select(&influxqb.Wildcard{}).From("MyMeasurement"),
		"SELECT * FROM MyMeasurement",
		false,
	},
	{
		"Select statement with unicode names",
		influxqb.NewSelectBuilder().Select("FieldA↓").From("MyMeasurement"),
		"SELECT \"FieldA↓\" FROM MyMeasurement",
		false,
	},
	{
		"Select statement with quotes in names",
		influxqb.NewSelectBuilder().Select("Fiel'dA").From("MyMea'surement"),
		"SELECT \"Fiel'dA\" FROM \"MyMea'surement\"",
		false,
	},
	{
		"Select statement with double quotes in names",
		influxqb.NewSelectBuilder().Select("Fiel\"dA").From("MyMea\"surement"),
		"SELECT \"Fiel\\\"dA\" FROM \"MyMea\\\"surement\"",
		false,
	},
	{
		"Select statement with regex as field",
		influxqb.NewSelectBuilder().Select(regexp.MustCompile(".*")).From("MyMeasurement"),
		"SELECT /.*/ FROM MyMeasurement",
		false,
	},
	{
		"Select statement with regex as MeasurementName",
		influxqb.NewSelectBuilder().Select("FieldA").FromRegex(regexp.MustCompile(".*")),
		"SELECT FieldA FROM /.*/",
		false,
	},
	{
		"Select function and string argument",
		influxqb.NewSelectBuilder().Select(&influxqb.Function{Name: "MEAN", Args: []interface{}{"FieldA"}}).From("MyMeasurement"),
		"SELECT MEAN('FieldA') FROM MyMeasurement",
		false,
	},
	{
		"Select function with Field, time, number arguments",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Function{Name: "MEAN", Args: []interface{}{&influxqb.Field{Name: "COL"}, 12, time.Date(2015, 8, 18, 0, 0, 0, 0, time.UTC)}}).
			From("MyMeasurement"),
		"SELECT MEAN(COL, 12, '2015-08-18T00:00:00Z') FROM MyMeasurement",
		false,
	},
	{
		"Select function with complex Field, time, number arguments",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Function{Name: "MEAN", Args: []interface{}{&influxqb.Field{Name: "CO'L"}, 12, time.Date(2015, 8, 18, 0, 0, 0, 0, time.UTC)}}).
			From("MyMeasurement"),
		"SELECT MEAN(\"CO'L\", 12, '2015-08-18T00:00:00Z') FROM MyMeasurement",
		false,
	},
	{
		"Select function with no fill from str",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Function{Name: "MEAN", Args: []interface{}{&influxqb.Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(influxqb.FillFromStr("none")),
		"SELECT MEAN(COL) FROM MyMeasurement fill(none)",
		false,
	},
	{
		"Select function with no fill",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Function{Name: "MEAN", Args: []interface{}{&influxqb.Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(&influxqb.FillNoFill{}),
		"SELECT MEAN(COL) FROM MyMeasurement fill(none)",
		false,
	},
	{
		"Select function with fill from str number",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Function{Name: "MEAN", Args: []interface{}{&influxqb.Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(influxqb.FillFromStr("123")),
		"SELECT MEAN(COL) FROM MyMeasurement fill(123)",
		false,
	},
	{
		"Select function with fill",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Function{Name: "MEAN", Args: []interface{}{&influxqb.Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(&influxqb.FillNumber{Value: 123}),
		"SELECT MEAN(COL) FROM MyMeasurement fill(123)",
		false,
	},
	{
		"Select function with fill from str previous",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Function{Name: "MEAN", Args: []interface{}{&influxqb.Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(influxqb.FillFromStr("previous")),
		"SELECT MEAN(COL) FROM MyMeasurement fill(previous)",
		false,
	},
	{
		"Select function with fill previous",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Function{Name: "MEAN", Args: []interface{}{&influxqb.Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(influxqb.FillPrevious{}),
		"SELECT MEAN(COL) FROM MyMeasurement fill(previous)",
		false,
	},
	{
		"Select function with fill from Str linear",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Function{Name: "MEAN", Args: []interface{}{&influxqb.Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(influxqb.FillFromStr("linear")),
		"SELECT MEAN(COL) FROM MyMeasurement fill(linear)",
		false,
	},
	{
		"Select function with fill linear",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Function{Name: "MEAN", Args: []interface{}{&influxqb.Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(influxqb.FillLinear{}),
		"SELECT MEAN(COL) FROM MyMeasurement fill(linear)",
		false,
	},
	{
		"Select function with fill null str",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Function{Name: "MEAN", Args: []interface{}{&influxqb.Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(influxqb.FillFromStr("null")),
		"SELECT MEAN(COL) FROM MyMeasurement",
		false,
	},
	{
		"Select function with fill null",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Function{Name: "MEAN", Args: []interface{}{&influxqb.Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(influxqb.FillNull{}),
		"SELECT MEAN(COL) FROM MyMeasurement",
		false,
	},
	{
		"Select function with fill unknown value str",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Function{Name: "MEAN", Args: []interface{}{&influxqb.Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(influxqb.FillFromStr("nonexistant")),
		"SELECT MEAN(COL) FROM MyMeasurement fill(none)",
		false,
	},
	{
		"Select function groupby time sampling and Field fill 123",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Function{Name: "MEAN", Args: []interface{}{&influxqb.Field{Name: "COL"}}},
		).
			From("MyMeasurement").
			GroupBy(&influxqb.Field{Name: "COLA"}, &influxqb.TimeSampling{Interval: time.Hour}).
			Fill(&influxqb.FillNumber{Value: 123}),
		"SELECT MEAN(COL) FROM MyMeasurement GROUP BY COLA, time(1h) fill(123)",
		false,
	},
	{
		"Select * with offset",
		influxqb.NewSelectBuilder().Select(&influxqb.Wildcard{}).
			From("MyMeasurement").
			Offset(125),
		"SELECT * FROM MyMeasurement OFFSET 125",
		false,
	},
	{
		"Select * with offset",
		influxqb.NewSelectBuilder().Select(&influxqb.Wildcard{}).
			From("MyMeasurement").
			Limit(125),
		"SELECT * FROM MyMeasurement LIMIT 125",
		false,
	},
	{
		"Select * with series offset",
		influxqb.NewSelectBuilder().Select(&influxqb.Wildcard{}).
			From("MyMeasurement").
			SeriesOffset(125),
		"SELECT * FROM MyMeasurement SOFFSET 125",
		false,
	},
	{
		"Select * with series offset",
		influxqb.NewSelectBuilder().Select(&influxqb.Wildcard{}).
			From("MyMeasurement").
			SeriesLimit(125),
		"SELECT * FROM MyMeasurement SLIMIT 125",
		false,
	},
	{
		"Select Math from measurement",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Math{Expr: []interface{}{&influxqb.Field{Name: "FieldC"}, influxql.ADD, 51}},
		).
			From("MyMeasurement").
			SeriesLimit(125),
		"SELECT (FieldC + 51) FROM MyMeasurement SLIMIT 125",
		false,
	},
	{
		"Select and order by time asc",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", influxqb.ASC),
		"SELECT * FROM MyMeasurement ORDER BY time ASC",
		false,
	},
	{
		"Select and order by time DESC",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", influxqb.DESC),
		"SELECT * FROM MyMeasurement ORDER BY time DESC",
		false,
	},
	{
		"Select and order by time DESC with timeZone",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", influxqb.DESC).
			WithTimeZone(tz),
		"SELECT * FROM MyMeasurement ORDER BY time DESC TZ('Europe/Paris')",
		false,
	},
	{
		"Select and order by time DESC with timeZone str",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", influxqb.DESC).
			WithTimeZone("Europe/Paris"),
		"SELECT * FROM MyMeasurement ORDER BY time DESC TZ('Europe/Paris')",
		false,
	},
	{
		"Select * where field is int and time less than time.Time",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", influxqb.DESC).
			Where(
				&influxqb.Math{Expr: []interface{}{
					&influxqb.Field{Name: "toto"}, influxql.EQ, 56,
					influxql.AND, &influxqb.Field{Name: "time"}, influxql.LT, time.Date(2020, 05, 16, 0, 0, 0, 153, time.UTC)},
				}),
		"SELECT * FROM MyMeasurement WHERE (toto = 56 AND time < '2020-05-16T00:00:00.000000153Z') ORDER BY time DESC",
		false,
	},
	{
		"Select * where complex",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", influxqb.DESC).
			Where(
				&influxqb.Math{Expr: []interface{}{
					&influxqb.Field{Name: "toto"}, influxql.EQ, 56,
					influxql.AND, &influxqb.Field{Name: "time"}, influxql.LT, time.Date(2020, 05, 16, 0, 0, 0, 153, time.UTC),
					influxql.OR,
					influxql.LPAREN,
					"tutututu", influxql.EQ, 12,
					influxql.AND,
					"aaa", influxql.EQ, "A",
					influxql.RPAREN,
					influxql.AND,
					influxql.LPAREN,
					&influxqb.Field{Name: "value"}, influxql.GTE, int32(323),
					influxql.AND, &influxqb.Field{Name: "computer"}, influxql.EQ, "toto",
					influxql.AND, &influxqb.Field{Name: "ptio"}, influxql.GT,
					influxql.LPAREN,
					15, influxql.ADD, 35.3,
					influxql.RPAREN,
					influxql.RPAREN}}),
		"SELECT * FROM MyMeasurement WHERE (toto = 56 AND time < '2020-05-16T00:00:00.000000153Z' OR ('tutututu' = 12 AND 'aaa' = 'A') AND (value >= 323 AND computer = 'toto' AND ptio > (15 + 35.300))) ORDER BY time DESC",
		false,
	},
	{
		"Select * where complex with Parenthesis object",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", influxqb.DESC).
			Where(
				&influxqb.Math{Expr: []interface{}{
					&influxqb.Field{Name: "toto"}, influxql.EQ, 56,
					influxql.AND, &influxqb.Field{Name: "time"}, influxql.LT, time.Date(2020, 05, 16, 0, 0, 0, 153, time.UTC),
					influxql.OR,
					&influxqb.Parenthesis{
						Expr: []interface{}{
							"tutututu", influxql.EQ, 12,
							influxql.AND,
							"aaa", influxql.EQ, "A",
						}},
					influxql.AND,
					&influxqb.Parenthesis{
						Expr: []interface{}{
							&influxqb.Field{Name: "value"}, influxql.GTE, int32(323),
							influxql.AND, &influxqb.Field{Name: "computer"}, influxql.EQ, "toto",
							influxql.AND, &influxqb.Field{Name: "ptio"}, influxql.GT, &influxqb.Parenthesis{Expr: []interface{}{uint16(15), influxql.ADD, 35.3}},
						}},
				}}),
		"SELECT * FROM MyMeasurement WHERE (toto = 56 AND time < '2020-05-16T00:00:00.000000153Z' OR ('tutututu' = 12 AND 'aaa' = 'A') AND (value >= 323 AND computer = 'toto' AND ptio > (15 + 35.300))) ORDER BY time DESC",
		false,
	},
	{
		"Select * where complex with Parenthesis object and Methods",
		influxqb.NewSelectBuilder().Select(
			&influxqb.Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", influxqb.DESC).
			Where(
				&influxqb.Math{Expr: []interface{}{
					influxqb.Eq(&influxqb.Field{Name: "toto"}, 56),
					influxql.AND,
					influxqb.LessThan(&influxqb.Field{Name: "time"}, time.Date(2020, 05, 16, 0, 0, 0, 153, time.UTC)),
					influxql.OR,
					&influxqb.Parenthesis{
						Expr: []interface{}{
							influxqb.And(
								influxqb.Eq("tutututu", 12),
								influxqb.Eq("aaa", "A")),
						}},
					influxql.AND,
					&influxqb.Parenthesis{
						Expr: []interface{}{
							influxqb.GreaterThanEq(&influxqb.Field{Name: "value"}, int32(323)),
							influxql.AND,
							influxqb.Eq(&influxqb.Field{Name: "computer"}, "toto"),
							influxql.AND,
							influxqb.GreaterThan(&influxqb.Field{Name: "ptio"}, &influxqb.Parenthesis{Expr: []interface{}{influxqb.Add(int8(15), 35.3)}}),
						}},
				}}),
		"SELECT * FROM MyMeasurement WHERE (toto = 56 AND time < '2020-05-16T00:00:00.000000153Z' OR ('tutututu' = 12 AND 'aaa' = 'A') AND (value >= 323 AND computer = 'toto' AND ptio > (15 + 35.300))) ORDER BY time DESC",
		false,
	},
	{
		"Select * where is string",
		influxqb.NewSelectBuilder().Select(&influxqb.Wildcard{}).From("MyMeasurement").OrderBy("time", influxqb.DESC).
			Where(&influxqb.MathExpr{Expr: "toto = 56 AND time < '2020-05-16T00:00:00.000000153Z' OR ('tutututu' = 12 AND 'aaa' = 'A') AND (value >= 323 AND computer = 'toto' AND ptio > (15 + 35.300))"}),
		"SELECT * FROM MyMeasurement WHERE (toto = 56 AND time < '2020-05-16T00:00:00.000000153Z' OR ('tutututu' = 12 AND 'aaa' = 'A') AND (value >= 323 AND computer = 'toto' AND ptio > (15 + 35.300))) ORDER BY time DESC",
		false,
	},
	{
		"Select Math expr",
		influxqb.NewSelectBuilder().Select(&influxqb.MathExpr{Expr: "FieldC + 12", Alias: "math"}).From("MyMeasurement"),
		"SELECT (FieldC + 12) AS math FROM MyMeasurement",
		false,
	},
	{
		"Select Math ",
		influxqb.NewSelectBuilder().Select(&influxqb.Math{Expr: []interface{}{
			influxqb.Add(&influxqb.Field{Name: "FieldC"}, 12)},
			Alias: "math"}).From("MyMeasurement"),
		"SELECT (FieldC + 12) AS math FROM MyMeasurement",
		false,
	},
	{
		"Select Add Select ",
		influxqb.NewSelectBuilder().Select("field1").AddSelect("field2").From("MyMeasurement"),
		"SELECT field1, field2 FROM MyMeasurement",
		false,
	},
	{
		"Select Add From ",
		influxqb.NewSelectBuilder().Select("field1").From("MyMeasurement").AddFrom("From2"),
		"SELECT field1 FROM MyMeasurement, From2",
		false,
	},
	{
		"Select Add From Regex ",
		influxqb.NewSelectBuilder().Select("field1").From("MyMeasurement").AddFromRegex(regexp.MustCompile(".*")),
		"SELECT field1 FROM MyMeasurement, /.*/",
		false,
	},
	{
		"Select Add GroupBy",
		influxqb.NewSelectBuilder().Select("field1").From("MyMeasurement").GroupBy(&influxqb.Field{Name: "f1"}).
			AddGroupBy(&influxqb.TimeSampling{Interval: time.Hour}),
		"SELECT field1 FROM MyMeasurement GROUP BY f1, time(1h)",
		false,
	},
	{
		"Select Subquery",
		influxqb.NewSelectBuilder().Select(&influxqb.Wildcard{}).
			FromSubQuery(influxqb.NewSelectBuilder().Select("F2", "F3").From("Table")),
		"SELECT * FROM (SELECT F2, F3 FROM Table)",
		false,
	},
	{
		"Select add Subquery",
		influxqb.NewSelectBuilder().Select(&influxqb.Wildcard{}).From("table2").
			AddFromSubQuery(influxqb.NewSelectBuilder().Select("F2", "F3").From("Table")),
		"SELECT * FROM table2, (SELECT F2, F3 FROM Table)",
		false,
	},
	{
		"Select WHERE or ",
		influxqb.NewSelectBuilder().Select(&influxqb.Wildcard{}).From("table2").
			Where(&influxqb.Math{Expr: []interface{}{influxqb.Or(
				influxqb.Eq(influxqb.Field{Name: "A"}, int16(165)),
				influxqb.LessThanEq(influxqb.Field{Name: "time"}, time.Date(1970, 01, 01, 0, 0, 0, 0, time.UTC)),
			)}}),
		"SELECT * FROM table2 WHERE (A = 165 OR time <= '1970-01-01T00:00:00Z')",
		false,
	},
	{
		"Select Divide, Subtract, multiply, modulus  Where noteq  ",
		influxqb.NewSelectBuilder().
			Select(&influxqb.Math{Expr: []interface{}{influxqb.Add(influxqb.Field{Name: "A"}, int64(32))}}).
			AddSelect(&influxqb.Math{Expr: []interface{}{influxqb.Divide(influxqb.Field{Name: "TU"}, float32(3.2))}}).
			AddSelect(&influxqb.Math{Expr: []interface{}{influxqb.Subtract(influxqb.Field{Name: "MINUS"}, uint(32))}}).
			AddSelect(&influxqb.Math{Expr: []interface{}{influxqb.Multiply(influxqb.Field{Name: "MINUS"}, uint64(1))}}).
			AddSelect(&influxqb.Math{Expr: []interface{}{influxqb.Modulus(influxqb.Field{Name: "MINUS"}, uint8(1))}}).
			From("table2").
			Where(&influxqb.Math{Expr: []interface{}{influxqb.Or(
				influxqb.NotEq(influxqb.Field{Name: "A"}, int16(165)),
				influxqb.LessThanEq(influxqb.Field{Name: "time"}, time.Date(1970, 01, 01, 0, 0, 0, 0, time.UTC)),
			)}}),
		"SELECT (A + 32), (TU / 3.200), (MINUS - 32), (MINUS * 1), (MINUS % 1) FROM table2 WHERE (A != 165 OR time <= '1970-01-01T00:00:00Z')",
		false,
	},
	{
		"Select Divide, Subtract, multiply, modulus  Where noteq No math object ",
		influxqb.NewSelectBuilder().
			Select(influxqb.Add(influxqb.Field{Name: "A"}, int64(32))).
			AddSelect(influxqb.Divide(influxqb.Field{Name: "TU"}, float32(3.2))).
			AddSelect(influxqb.Subtract(influxqb.Field{Name: "MINUS"}, uint(32))).
			AddSelect(influxqb.Multiply(influxqb.Field{Name: "MINUS"}, uint64(1))).
			AddSelect(influxqb.Modulus(influxqb.Field{Name: "MINUS"}, uint8(1))).
			From("table2").
			Where(
				influxqb.Or(
					influxqb.NotEq(influxqb.Field{Name: "A"}, int16(165)),
					influxqb.LessThanEq(influxqb.Field{Name: "time"}, time.Date(1970, 01, 01, 0, 0, 0, 0, time.UTC)),
				)),
		"SELECT (A + 32), (TU / 3.200), (MINUS - 32), (MINUS * 1), (MINUS % 1) FROM table2 WHERE (A != 165 OR time <= '1970-01-01T00:00:00Z')",
		false,
	},
	{
		"Select into ",
		influxqb.NewSelectBuilder().
			Select(&influxqb.Field{Name: "A"}).
			From("table2").
			Into(&influxql.Measurement{
				Database:        "MyDB",
				RetentionPolicy: "RP",
				Name:            "Measurement",
				IsTarget:        true,
			}),
		"SELECT A INTO MyDB.RP.\"Measurement\" FROM table2",
		false,
	},
}

func TestSelect(t *testing.T) {
	for i, sample := range testSamples {
		s, err := sample.b.Build()

		fmt.Print("Test ", i, ": ", sample.d)

		if sample.e {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		assert.Equal(t, sample.s, s)

		fmt.Println("   [OK]")
	}
}
