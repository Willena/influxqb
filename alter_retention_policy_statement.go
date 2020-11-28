package influxqb

import (
	"github.com/influxdata/influxql"
	"time"
)

/*

	// Duration of the Shard.
	ShardGroupDuration *time.Duration
*/

type AlterRetentionPolicyBuilder struct {
	alterStm *influxql.AlterRetentionPolicyStatement
}

func (b *AlterRetentionPolicyBuilder) WithDatabase(database string) *AlterRetentionPolicyBuilder {
	b.alterStm.Database = database
	return b
}

func (b *AlterRetentionPolicyBuilder) WithPolicyName(policyName string) *AlterRetentionPolicyBuilder {
	b.alterStm.Name = policyName
	return b
}

func (b *AlterRetentionPolicyBuilder) WithDuration(duration *time.Duration) *AlterRetentionPolicyBuilder {
	b.alterStm.Duration = duration
	return b
}
func (b *AlterRetentionPolicyBuilder) WithDurationString(durationStr string) *AlterRetentionPolicyBuilder {
	duration := new(time.Duration)
	*duration, _ = time.ParseDuration(durationStr)
	b.WithDuration(duration)
	return b
}

func (b *AlterRetentionPolicyBuilder) WithShardDuration(duration *time.Duration) *AlterRetentionPolicyBuilder {
	b.alterStm.ShardGroupDuration = duration
	return b
}

func (b *AlterRetentionPolicyBuilder) WithShardDurationString(durationStr string) *AlterRetentionPolicyBuilder {
	duration := new(time.Duration)
	*duration, _ = time.ParseDuration(durationStr)
	b.WithShardDuration(duration)
	return b
}

func (b *AlterRetentionPolicyBuilder) WithReplicationFactor(replication int) *AlterRetentionPolicyBuilder {
	n := new(int)
	*n = replication
	b.alterStm.Replication = n
	return b
}

func (b *AlterRetentionPolicyBuilder) SetAsDefault() *AlterRetentionPolicyBuilder {
	b.alterStm.Default = true
	return b
}

func (b *AlterRetentionPolicyBuilder) Build() (string, error) {
	return b.alterStm.String(), nil
}
