# Influxqb - A simple Query builder for InfluxDB

## Intro

Building strongly typed and secure InfluxQL queries from string is not always an easy task
You have to take care of the sanitization, keep the distinction between function, identifier, numbers and string literals
to build a valid query. 

The influxQL parser contains all the required types to manually build a query from scratch. More importantly it also contains 
`String()` method on each type and statement. The parser then take care of escaping character, putting quotes or not, ...

This go package is built on top of the influxql parser and offers a more simple way to create queries.  

## Things that are working

* select_stmt 
* alter_retention_policy_stmt
* create_continuous_query_stmt
* create_retention_policy_stmt
* create_database_stmt
* create_subscription_stmt
* create_user_stmt
* delete_stmt
* drop_continuous_query_stmt 
* drop_database_stmt 
* drop_measurement_stmt 
* drop_retention_policy_stmt 
* drop_series_stmt 


## Todo


* drop_shard_stmt 
* drop_subscription_stmt 
* drop_user_stmt 
* explain_stmt 
* explain_analyze_stmt 
* grant_stmt 
* kill_query_statement 
* revoke_stmt 
* show_continuous_queries_stmt 
* show_databases_stmt 
* show_diagnostics_stmt 
* show_field_key_cardinality_stmt 
* show_field_keys_stmt 
* show_grants_stmt 
* show_measurement_cardinality_stmt 
* show_measurement_exact_cardinality_stmt 
* show_measurements_stmt 
* show_queries_stmt 
* show_retention_policies_stmt 
* show_series_cardinality_stmt 
* show_series_exact_cardinality_stmt 
* show_series_stmt 
* show_shard_groups_stmt 
* show_shards_stmt 
* show_stats_stmt 
* show_subscriptions_stmt 
* show_tag_key_cardinality_stmt 
* show_tag_key_exact_cardinality_stmt 
* show_tag_keys_stmt 
* show_tag_values_stmt 
* show_tag_values_cardinality_stmt 
* show_users_stmt 


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
