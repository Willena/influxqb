# Influxqb - A simple Query builder for InfluxDB

## Intro

Building strongly typed and secure InfluxQL queries from string is not always an easy task
You have to take care of the sanitization, keep the distinction between function, identifier, numbers and string literals
to build a valid query. 

The influxQL parser contains all the required types to manually build a query from scratch. More importantly it also contains 
`String()` method on each type and statement. The parser then take care of escaping character, putting quotes or not, ...

This go package is built on top of the influxql parser and offers a more simple way to create queries conveniently.  

## example 
```go
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
    &influxqb.Math{Expr: []interface{}{influxqb.Field{Name:"Tptp"}, influxql.EQ, "data", influxql.AND, "ooo", influxql.EQ, 16.55},
})
```