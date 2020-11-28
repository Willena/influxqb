package influxqb

import "github.com/influxdata/influxql"

type CreateDatabaseBuilder struct {
	dbStatement *influxql.CreateDatabaseStatement
}

func (b *CreateDatabaseBuilder) WithName(str string) *CreateDatabaseBuilder {
	b.dbStatement.Name = str
	return b
}

func (b *CreateDatabaseBuilder) WithRetentionPolicy(retentionPolicy *CreateRetentionPolicyBuilder) *CreateDatabaseBuilder {

	b.dbStatement.RetentionPolicyCreate = true
	b.dbStatement.RetentionPolicyName = retentionPolicy.ret.Name
	b.dbStatement.RetentionPolicyDuration = &retentionPolicy.ret.Duration
	b.dbStatement.RetentionPolicyReplication = &retentionPolicy.ret.Replication
	b.dbStatement.RetentionPolicyShardGroupDuration = retentionPolicy.ret.ShardGroupDuration

	return b
}

func (b *CreateDatabaseBuilder) Build() (string, error) {
	return b.dbStatement.String(), nil
}
