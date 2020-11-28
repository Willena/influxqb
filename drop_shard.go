package influxqb

import "github.com/influxdata/influxql"

type DropShardBuilder struct {
	dss *influxql.DropShardStatement
}

func (b *DropShardBuilder) WithShard(id uint64) *DropShardBuilder {
	b.dss.ID = id
	return b
}

func (b *DropShardBuilder) Build() (string, error) {
	return b.dss.String(), nil
}
