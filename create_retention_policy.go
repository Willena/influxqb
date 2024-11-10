package influxqb

import (
	"github.com/influxdata/influxql"
	"time"
)

type CreateRetentionPolicyBuilder struct {
	ret *influxql.CreateRetentionPolicyStatement
}

func (b *CreateRetentionPolicyBuilder) WithDatabase(database string) *CreateRetentionPolicyBuilder {
	b.ret.Database = database
	return b
}

func (b *CreateRetentionPolicyBuilder) WithPolicyName(policyName string) *CreateRetentionPolicyBuilder {
	b.ret.Name = policyName
	return b
}

func (b *CreateRetentionPolicyBuilder) WithDuration(duration time.Duration) *CreateRetentionPolicyBuilder {
	b.ret.Duration = duration
	return b
}
func (b *CreateRetentionPolicyBuilder) WithDurationString(durationStr string) *CreateRetentionPolicyBuilder {
	duration, _ := time.ParseDuration(durationStr)
	b.WithDuration(duration)
	return b
}

func (b *CreateRetentionPolicyBuilder) WithShardDuration(duration time.Duration) *CreateRetentionPolicyBuilder {
	b.ret.ShardGroupDuration = duration
	return b
}

func (b *CreateRetentionPolicyBuilder) WithShardDurationString(durationStr string) *CreateRetentionPolicyBuilder {
	duration, _ := time.ParseDuration(durationStr)
	b.WithShardDuration(duration)
	return b
}

func (b *CreateRetentionPolicyBuilder) WithReplicationFactor(replication int) *CreateRetentionPolicyBuilder {
	b.ret.Replication = replication
	return b
}

func (b *CreateRetentionPolicyBuilder) WithFutureLimit(futureLimit time.Duration) *CreateRetentionPolicyBuilder {
	b.ret.FutureWriteLimit = futureLimit
	return b
}

func (b *CreateRetentionPolicyBuilder) WithFutureLimitString(futureLimit string) *CreateRetentionPolicyBuilder {
	duration, _ := time.ParseDuration(futureLimit)
	b.WithFutureLimit(duration)
	return b
}

func (b *CreateRetentionPolicyBuilder) WithPastLimit(pastLimit time.Duration) *CreateRetentionPolicyBuilder {
	b.ret.PastWriteLimit = pastLimit
	return b
}

func (b *CreateRetentionPolicyBuilder) WithPastLimitString(pastLimitStr string) *CreateRetentionPolicyBuilder {
	duration, _ := time.ParseDuration(pastLimitStr)
	b.WithPastLimit(duration)
	return b
}

func (b *CreateRetentionPolicyBuilder) SetAsDefault() *CreateRetentionPolicyBuilder {
	b.ret.Default = true
	return b
}

func (b *CreateRetentionPolicyBuilder) Build() (string, error) {
	return b.ret.String(), nil
}
