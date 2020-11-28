package influxqb

import "github.com/influxdata/influxql"

type DropDatabaseBuilder struct {
	dcq *influxql.DropDatabaseStatement
}

func (b *DropDatabaseBuilder) WithDatabase(database string) *DropDatabaseBuilder {
	b.dcq.Name = database
	return b
}

func (b *DropDatabaseBuilder) Build() (string, error) {
	return b.dcq.String(), nil
}
