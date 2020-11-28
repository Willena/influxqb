package influxqb

import "github.com/influxdata/influxql"

type DropContinuousQueryBuilder struct {
	dcq *influxql.DropContinuousQueryStatement
}

func (b *DropContinuousQueryBuilder) WithName(name string) *DropContinuousQueryBuilder {
	b.dcq.Name = name
	return b
}
func (b *DropContinuousQueryBuilder) WithDatabase(database string) *DropContinuousQueryBuilder {
	b.dcq.Database = database
	return b
}

func (b *DropContinuousQueryBuilder) Build() (string, error) {
	return b.dcq.String(), nil
}
