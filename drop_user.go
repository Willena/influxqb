package influxqb

import "github.com/influxdata/influxql"

type DropUserBuilder struct {
	dss *influxql.DropUserStatement
}

func (b *DropUserBuilder) WithUsername(username string) *DropUserBuilder {
	b.dss.Name = username
	return b
}

func (b *DropUserBuilder) Build() (string, error) {
	return b.dss.String(), nil
}
