package influxqb

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/influxdata/influxql"
	"github.com/stretchr/testify/assert"
)

var tz, _ = time.LoadLocation("Europe/Paris")

var testSamples = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{
		"Simple select statement",
		NewSelectBuilder().Select("FieldA").From("MyMeasurement"),
		"SELECT FieldA FROM MyMeasurement",
		false,
	},
	{
		"Simple select * statement",
		NewSelectBuilder().Select(&Wildcard{}).From("MyMeasurement"),
		"SELECT * FROM MyMeasurement",
		false,
	},
	{
		"Select statement with unicode names",
		NewSelectBuilder().Select("FieldA↓").From("MyMeasurement"),
		"SELECT \"FieldA↓\" FROM MyMeasurement",
		false,
	},
	{
		"Select statement with quotes in names",
		NewSelectBuilder().Select("Fiel'dA").From("MyMea'surement"),
		"SELECT \"Fiel'dA\" FROM \"MyMea'surement\"",
		false,
	},
	{
		"Select statement with double quotes in names",
		NewSelectBuilder().Select("Fiel\"dA").From("MyMea\"surement"),
		"SELECT \"Fiel\\\"dA\" FROM \"MyMea\\\"surement\"",
		false,
	},
	{
		"Select statement with regex as field",
		NewSelectBuilder().Select(regexp.MustCompile(".*")).From("MyMeasurement"),
		"SELECT /.*/ FROM MyMeasurement",
		false,
	},
	{
		"Select statement with regex as MeasurementName",
		NewSelectBuilder().Select("FieldA").FromRegex(regexp.MustCompile(".*")),
		"SELECT FieldA FROM /.*/",
		false,
	},
	{
		"Select function and string argument",
		NewSelectBuilder().Select(NewFunction("MEAN").WithArgs("FieldA")).From("MyMeasurement"),
		"SELECT MEAN('FieldA') FROM MyMeasurement",
		false,
	},
	{
		"Select function with Field, time, number arguments",
		NewSelectBuilder().Select(
			&Function{Name: "MEAN", Args: []interface{}{&Field{Name: "COL"}, 12, time.Date(2015, 8, 18, 0, 0, 0, 0, time.UTC)}}).
			From("MyMeasurement"),
		"SELECT MEAN(COL, 12, '2015-08-18T00:00:00Z') FROM MyMeasurement",
		false,
	},
	{
		"Select function with complex Field, time, number arguments",
		NewSelectBuilder().Select(
			&Function{Name: "MEAN", Args: []interface{}{&Field{Name: "CO'L"}, 12, time.Date(2015, 8, 18, 0, 0, 0, 0, time.UTC)}}).
			From("MyMeasurement"),
		"SELECT MEAN(\"CO'L\", 12, '2015-08-18T00:00:00Z') FROM MyMeasurement",
		false,
	},
	{
		"Select function with no fill from str",
		NewSelectBuilder().Select(
			&Function{Name: "MEAN", Args: []interface{}{&Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(FillFromStr("none")),
		"SELECT MEAN(COL) FROM MyMeasurement fill(none)",
		false,
	},
	{
		"Select function with no fill",
		NewSelectBuilder().Select(
			&Function{Name: "MEAN", Args: []interface{}{&Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(&FillNoFill{}),
		"SELECT MEAN(COL) FROM MyMeasurement fill(none)",
		false,
	},
	{
		"Select function with fill from str number",
		NewSelectBuilder().Select(
			&Function{Name: "MEAN", Args: []interface{}{&Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(FillFromStr("123")),
		"SELECT MEAN(COL) FROM MyMeasurement fill(123)",
		false,
	},
	{
		"Select function with fill",
		NewSelectBuilder().Select(
			&Function{Name: "MEAN", Args: []interface{}{&Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(&FillNumber{Value: 123}),
		"SELECT MEAN(COL) FROM MyMeasurement fill(123)",
		false,
	},
	{
		"Select function with fill from str previous",
		NewSelectBuilder().Select(
			&Function{Name: "MEAN", Args: []interface{}{&Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(FillFromStr("previous")),
		"SELECT MEAN(COL) FROM MyMeasurement fill(previous)",
		false,
	},
	{
		"Select function with fill previous",
		NewSelectBuilder().Select(
			&Function{Name: "MEAN", Args: []interface{}{&Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(FillPrevious{}),
		"SELECT MEAN(COL) FROM MyMeasurement fill(previous)",
		false,
	},
	{
		"Select function with fill from Str linear",
		NewSelectBuilder().Select(
			&Function{Name: "MEAN", Args: []interface{}{&Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(FillFromStr("linear")),
		"SELECT MEAN(COL) FROM MyMeasurement fill(linear)",
		false,
	},
	{
		"Select function with fill linear",
		NewSelectBuilder().Select(
			&Function{Name: "MEAN", Args: []interface{}{&Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(FillLinear{}),
		"SELECT MEAN(COL) FROM MyMeasurement fill(linear)",
		false,
	},
	{
		"Select function with fill null str",
		NewSelectBuilder().Select(
			&Function{Name: "MEAN", Args: []interface{}{&Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(FillFromStr("null")),
		"SELECT MEAN(COL) FROM MyMeasurement",
		false,
	},
	{
		"Select function with fill null",
		NewSelectBuilder().Select(
			&Function{Name: "MEAN", Args: []interface{}{&Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(FillNull{}),
		"SELECT MEAN(COL) FROM MyMeasurement",
		false,
	},
	{
		"Select function with fill unknown value str",
		NewSelectBuilder().Select(
			&Function{Name: "MEAN", Args: []interface{}{&Field{Name: "COL"}}}).
			From("MyMeasurement").Fill(FillFromStr("nonexistant")),
		"SELECT MEAN(COL) FROM MyMeasurement fill(none)",
		false,
	},
	{
		"Select function groupby time sampling and Field fill 123",
		NewSelectBuilder().Select(
			&Function{Name: "MEAN", Args: []interface{}{&Field{Name: "COL"}}},
		).
			From("MyMeasurement").
			GroupBy(&Field{Name: "COLA"}, &TimeSampling{Interval: time.Hour}).
			Fill(&FillNumber{Value: 123}),
		"SELECT MEAN(COL) FROM MyMeasurement GROUP BY COLA, time(1h) fill(123)",
		false,
	},
	{
		"Select * with offset",
		NewSelectBuilder().Select(&Wildcard{}).
			From("MyMeasurement").
			Offset(125),
		"SELECT * FROM MyMeasurement OFFSET 125",
		false,
	},
	{
		"Select * with limit",
		NewSelectBuilder().Select(&Wildcard{}).
			From("MyMeasurement").
			Limit(125),
		"SELECT * FROM MyMeasurement LIMIT 125",
		false,
	},
	{
		"Select * with series offset",
		NewSelectBuilder().Select(&Wildcard{}).
			From("MyMeasurement").
			SeriesOffset(125),
		"SELECT * FROM MyMeasurement SOFFSET 125",
		false,
	},
	{
		"Select * with series offset",
		NewSelectBuilder().Select(&Wildcard{}).
			From("MyMeasurement").
			SeriesLimit(125),
		"SELECT * FROM MyMeasurement SLIMIT 125",
		false,
	},
	{
		"Select Math from measurement",
		NewSelectBuilder().Select(
			&Math{Expr: []interface{}{&Field{Name: "FieldC"}, influxql.ADD, 51}},
		).
			From("MyMeasurement").
			SeriesLimit(125),
		"SELECT (FieldC + 51) FROM MyMeasurement SLIMIT 125",
		false,
	},
	{
		"Select and order by time asc",
		NewSelectBuilder().Select(
			&Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", ASC),
		"SELECT * FROM MyMeasurement ORDER BY time ASC",
		false,
	},
	{
		"Select and order by time DESC",
		NewSelectBuilder().Select(
			&Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", DESC),
		"SELECT * FROM MyMeasurement ORDER BY time DESC",
		false,
	},
	{
		"Select and order by time DESC with timeZone",
		NewSelectBuilder().Select(
			&Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", DESC).
			WithTimeZone(tz),
		"SELECT * FROM MyMeasurement ORDER BY time DESC TZ('Europe/Paris')",
		false,
	},
	{
		"Select and order by time DESC with timeZone str",
		NewSelectBuilder().Select(
			&Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", DESC).
			WithTimeZone("Europe/Paris"),
		"SELECT * FROM MyMeasurement ORDER BY time DESC TZ('Europe/Paris')",
		false,
	},
	{
		"Select * where field is int and time less than time.Time",
		NewSelectBuilder().Select(
			&Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", DESC).
			Where(
				&Math{Expr: []interface{}{
					&Field{Name: "toto"}, influxql.EQ, 56,
					influxql.AND, &Field{Name: "time"}, influxql.LT, time.Date(2020, 05, 16, 0, 0, 0, 153, time.UTC)},
				}),
		"SELECT * FROM MyMeasurement WHERE (toto = 56 AND time < '2020-05-16T00:00:00.000000153Z') ORDER BY time DESC",
		false,
	},
	{
		"Select * where complex",
		NewSelectBuilder().Select(
			&Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", DESC).
			Where(
				&Math{Expr: []interface{}{
					&Field{Name: "toto"}, influxql.EQ, 56,
					influxql.AND, &Field{Name: "time"}, influxql.LT, time.Date(2020, 05, 16, 0, 0, 0, 153, time.UTC),
					influxql.OR,
					influxql.LPAREN,
					"tutututu", influxql.EQ, 12,
					influxql.AND,
					"aaa", influxql.EQ, "A",
					influxql.RPAREN,
					influxql.AND,
					influxql.LPAREN,
					&Field{Name: "value"}, influxql.GTE, int32(323),
					influxql.AND, &Field{Name: "computer"}, influxql.EQ, "toto",
					influxql.AND, &Field{Name: "ptio"}, influxql.GT,
					influxql.LPAREN,
					15, influxql.ADD, 35.3,
					influxql.RPAREN,
					influxql.RPAREN}}),
		"SELECT * FROM MyMeasurement WHERE (toto = 56 AND time < '2020-05-16T00:00:00.000000153Z' OR ('tutututu' = 12 AND 'aaa' = 'A') AND (value >= 323 AND computer = 'toto' AND ptio > (15 + 35.3))) ORDER BY time DESC",
		false,
	},
	{
		"Select * where complex with Parenthesis object",
		NewSelectBuilder().Select(
			&Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", DESC).
			Where(
				&Math{Expr: []interface{}{
					&Field{Name: "toto"}, influxql.EQ, 56,
					influxql.AND, &Field{Name: "time"}, influxql.LT, time.Date(2020, 05, 16, 0, 0, 0, 153, time.UTC),
					influxql.OR,
					&Parenthesis{
						Expr: []interface{}{
							"tutututu", influxql.EQ, 12,
							influxql.AND,
							"aaa", influxql.EQ, "A",
						}},
					influxql.AND,
					&Parenthesis{
						Expr: []interface{}{
							&Field{Name: "value"}, influxql.GTE, int32(323),
							influxql.AND, &Field{Name: "computer"}, influxql.EQ, "toto",
							influxql.AND, &Field{Name: "ptio"}, influxql.GT, &Parenthesis{Expr: []interface{}{uint16(15), influxql.ADD, 35.3}},
						}},
				}}),
		"SELECT * FROM MyMeasurement WHERE (toto = 56 AND time < '2020-05-16T00:00:00.000000153Z' OR ('tutututu' = 12 AND 'aaa' = 'A') AND (value >= 323 AND computer = 'toto' AND ptio > (15 + 35.3))) ORDER BY time DESC",
		false,
	},
	{
		"Select * where complex with Parenthesis object and Methods",
		NewSelectBuilder().Select(
			&Wildcard{},
		).
			From("MyMeasurement").
			OrderBy("time", DESC).
			Where(
				&Math{Expr: []interface{}{
					Eq(&Field{Name: "toto"}, 56),
					influxql.AND,
					LessThan(&Field{Name: "time"}, time.Date(2020, 05, 16, 0, 0, 0, 153, time.UTC)),
					influxql.OR,
					&Parenthesis{
						Expr: []interface{}{
							And(
								Eq("tutututu", 12),
								Eq("aaa", "A")),
						}},
					influxql.AND,
					&Parenthesis{
						Expr: []interface{}{
							GreaterThanEq(&Field{Name: "value"}, int32(323)),
							influxql.AND,
							Eq(&Field{Name: "computer"}, "toto"),
							influxql.AND,
							GreaterThan(&Field{Name: "ptio"}, &Parenthesis{Expr: []interface{}{Add(int8(15), 35.3)}}),
						}},
				}}),
		"SELECT * FROM MyMeasurement WHERE (toto = 56 AND time < '2020-05-16T00:00:00.000000153Z' OR ('tutututu' = 12 AND 'aaa' = 'A') AND (value >= 323 AND computer = 'toto' AND ptio > (15 + 35.3))) ORDER BY time DESC",
		false,
	},
	{
		"Select * where is string",
		NewSelectBuilder().Select(NewWildcardField()).From("MyMeasurement").OrderBy("time", DESC).
			Where(&MathExpr{Expr: "toto = 56 AND time < '2020-05-16T00:00:00.000000153Z' OR ('tutututu' = 12 AND 'aaa' = 'A') AND (value >= 323 AND computer = 'toto' AND ptio > (15 + 35.3))"}),
		"SELECT * FROM MyMeasurement WHERE (toto = 56 AND time < '2020-05-16T00:00:00.000000153Z' OR ('tutututu' = 12 AND 'aaa' = 'A') AND (value >= 323 AND computer = 'toto' AND ptio > (15 + 35.3))) ORDER BY time DESC",
		false,
	},
	{
		"Select Math expr",
		NewSelectBuilder().Select(&MathExpr{Expr: "FieldC + 12", Alias: "math"}).From("MyMeasurement"),
		"SELECT (FieldC + 12) AS math FROM MyMeasurement",
		false,
	},
	{
		"Select Math ",
		NewSelectBuilder().Select(&Math{Expr: []interface{}{
			Add(&Field{Name: "FieldC"}, 12)},
			Alias: "math"}).From("MyMeasurement"),
		"SELECT (FieldC + 12) AS math FROM MyMeasurement",
		false,
	},
	{
		"Select Add Select ",
		NewSelectBuilder().Select("field1").AddSelect("field2").From("MyMeasurement"),
		"SELECT field1, field2 FROM MyMeasurement",
		false,
	},
	{
		"Select Add From ",
		NewSelectBuilder().Select("field1").From("MyMeasurement").AddFrom("From2"),
		"SELECT field1 FROM MyMeasurement, From2",
		false,
	},
	{
		"Select Add From Regex ",
		NewSelectBuilder().Select("field1").From("MyMeasurement").AddFromRegex(regexp.MustCompile(".*")),
		"SELECT field1 FROM MyMeasurement, /.*/",
		false,
	},
	{
		"Select Add GroupBy",
		NewSelectBuilder().Select("field1").From("MyMeasurement").GroupBy(&Field{Name: "f1"}).
			AddGroupBy(NewTimeSampling(time.Hour)),
		"SELECT field1 FROM MyMeasurement GROUP BY f1, time(1h)",
		false,
	},
	{
		"Select Subquery",
		NewSelectBuilder().Select(&Wildcard{}).
			FromSubQuery(NewSelectBuilder().Select("F2", "F3").From("Table")),
		"SELECT * FROM (SELECT F2, F3 FROM Table)",
		false,
	},
	{
		"Select add Subquery",
		NewSelectBuilder().Select(&Wildcard{}).From("table2").
			AddFromSubQuery(NewSelectBuilder().Select("F2", "F3").From("Table")),
		"SELECT * FROM table2, (SELECT F2, F3 FROM Table)",
		false,
	},
	{
		"Select WHERE or ",
		NewSelectBuilder().Select(&Wildcard{}).From("table2").
			Where(&Math{Expr: []interface{}{Or(
				Eq(Field{Name: "A"}, int16(165)),
				LessThanEq(Field{Name: "time"}, time.Date(1970, 01, 01, 0, 0, 0, 0, time.UTC)),
			)}}),
		"SELECT * FROM table2 WHERE (A = 165 OR time <= '1970-01-01T00:00:00Z')",
		false,
	},
	{
		"Select Divide, Subtract, multiply, modulus  Where noteq  ",
		NewSelectBuilder().
			Select(&Math{Expr: []interface{}{Add(Field{Name: "A"}, int64(32))}}).
			AddSelect(&Math{Expr: []interface{}{Divide(Field{Name: "TU"}, 3.3)}}).
			AddSelect(&Math{Expr: []interface{}{Subtract(Field{Name: "MINUS"}, uint(32))}}).
			AddSelect(&Math{Expr: []interface{}{Multiply(Field{Name: "MINUS"}, uint64(1))}}).
			AddSelect(&Math{Expr: []interface{}{Modulus(Field{Name: "MINUS"}, uint8(1))}}).
			From("table2").
			Where(NewMath().WithExpr(Or(
				NotEq(Field{Name: "A"}, int16(165)),
				LessThanEq(Field{Name: "time"}, time.Date(1970, 01, 01, 0, 0, 0, 0, time.UTC)),
			))),
		"SELECT (A + 32), (TU / 3.3), (MINUS - 32), (MINUS * 1), (MINUS % 1) FROM table2 WHERE (A != 165 OR time <= '1970-01-01T00:00:00Z')",
		false,
	},
	{
		"Select Divide, Subtract, multiply, modulus  Where noteq No math object ",
		NewSelectBuilder().
			Select(Add(Field{Name: "A"}, int64(32))).
			AddSelect(Divide(Field{Name: "TU"}, 3.3)).
			AddSelect(Subtract(Field{Name: "MINUS"}, uint(32))).
			AddSelect(Multiply(Field{Name: "MINUS"}, uint64(1))).
			AddSelect(Modulus(Field{Name: "MINUS"}, uint8(1))).
			From("table2").
			Where(
				Or(
					NotEq(NewField("A"), int16(165)),
					LessThanEq(Field{Name: "time"}, time.Date(1970, 01, 01, 0, 0, 0, 0, time.UTC)),
				)),
		"SELECT (A + 32), (TU / 3.3), (MINUS - 32), (MINUS * 1), (MINUS % 1) FROM table2 WHERE (A != 165 OR time <= '1970-01-01T00:00:00Z')",
		false,
	},
	{
		"Select into Nice functions ",
		NewSelectBuilder().
			Select(&Field{Name: "A"}).
			From("table2").
			Into(NewMeasurement().WithDatabase("MyDB").WithPolicy("RP").Name("Measurement")),
		"SELECT A INTO MyDB.RP.\"Measurement\" FROM table2",
		false,
	},
	{
		"Select from db.rp.name",
		NewSelectBuilder().Select(NewWildcardField()).FromMeasurements(NewMeasurement().Name("metricName").WithDatabase("db").WithPolicy("polName")),
		"SELECT * FROM db.polName.metricName",
		false,
	},
}

func TestSelectBuilder(t *testing.T) {
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
