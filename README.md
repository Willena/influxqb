# Influxqb - A simple Query builder for InfluxDB

## Intro

Building strongly typed and secure InfluxQL queries from string is not always an easy task
You have to take care of the sanitization, keep the distinction between function, identifier, numbers and string literals
to build a valid query. 

The influxQL parser contains all the required types to manually build a query from scratch. More importantly it also contains 
`String()` method on each type and statement. The parser then take care of escaping character, putting quotes or not, ...

This go package is built on top of the influxql parser and offers a more simple way to create queries.  

## Things that are working

* Select Statements 



## example 
```go
package main 

import (
	"fmt"
	"github.com/influxdata/influxql"
	"github.com/willena/influxqb"
	"time"
)

func main() {

	builder := influxqb.NewSelectBuilder()
	builder.Select(
		&influxqb.Function{Name: "MEAN", Args: []interface{}{"colomn", time.Now(), 45.36, time.Hour}},
		&influxqb.Field{Name: "MyField"},
	)
	builder.From("XTC_OLD'sk")
	builder.GroupBy(
		&influxqb.Field{Name: "GroupByField"},
		&influxqb.TimeSampling{Interval: time.Hour},
	)
	builder.Fill(45)
	builder.Limit(250)
	builder.Offset(15)
	builder.SeriesLimit(2)
	builder.SeriesOffset(8)
	builder.Where(
		influxqb.And(
			influxqb.Eq(influxqb.Field{Name: "Tptp"}, "data"),
			influxqb.Eq("ooo", 16.55)),
		)

	fmt.Println(builder.Build())
}
```

## Licence
   
   InfluxQb Go package
   Copyright 2020 Guillaume VILLENA aka "Willena"
   
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
   
      http://www.apache.org/licenses/LICENSE-2.0
   
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
