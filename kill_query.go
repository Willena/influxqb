package influxqb

import "github.com/influxdata/influxql"

type KillQueryBuilder struct {
	dss *influxql.KillQueryStatement
}

func (b *KillQueryBuilder) WithQueryId(id uint64) *KillQueryBuilder {
	b.dss.QueryID = id
	return b
}

func (b *KillQueryBuilder) OnHost(str string) *KillQueryBuilder {
	b.dss.Host = str
	return b
}

func (b *KillQueryBuilder) Build() (string, error) {
	return b.dss.String(), nil
}
