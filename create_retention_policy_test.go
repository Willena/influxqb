package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCreateRetentionPolicy = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{
		"Simple Retention policy",
		NewCreateRetentionPolicyBuilder().WithPolicyName("PolicyName").WithDatabase("database"),
		"CREATE RETENTION POLICY PolicyName ON \"database\" DURATION 0s REPLICATION 0",
		false,
	},
	{
		"Simple Retention policy set default",
		NewCreateRetentionPolicyBuilder().
			WithPolicyName("PolicyName").
			WithDatabase("database").
			SetAsDefault(),
		"CREATE RETENTION POLICY PolicyName ON \"database\" DURATION 0s REPLICATION 0 DEFAULT",
		false,
	},
	{
		"Simple Retention policy set duration",
		NewCreateRetentionPolicyBuilder().
			WithPolicyName("PolicyName").
			WithDatabase("database").
			WithDurationString("1h"),
		"CREATE RETENTION POLICY PolicyName ON \"database\" DURATION 1h REPLICATION 0",
		false,
	},
	{
		"Simple Retention policy set  shard duration",
		NewCreateRetentionPolicyBuilder().
			WithPolicyName("PolicyName").
			WithDatabase("database").
			WithShardDurationString("1h"),
		"CREATE RETENTION POLICY PolicyName ON \"database\" DURATION 0s REPLICATION 0 SHARD DURATION 1h",
		false,
	},
	{
		"Simple Retention policy set  shard duration and duration",
		NewCreateRetentionPolicyBuilder().
			WithPolicyName("PolicyName").
			WithDatabase("database").
			WithShardDurationString("1h").
			WithDurationString("3h"),
		"CREATE RETENTION POLICY PolicyName ON \"database\" DURATION 3h REPLICATION 0 SHARD DURATION 1h",
		false,
	},
	{
		"Simple Retention policy set  shard duration and duration plus replication factor",
		NewCreateRetentionPolicyBuilder().
			WithPolicyName("PolicyName").
			WithDatabase("database").
			WithShardDurationString("1h").
			WithDurationString("3h").
			WithReplicationFactor(36),
		"CREATE RETENTION POLICY PolicyName ON \"database\" DURATION 3h REPLICATION 36 SHARD DURATION 1h",
		false,
	},
	{
		"Simple Retention policy unicode names and quotes",
		NewCreateRetentionPolicyBuilder().
			WithPolicyName("PolicyName↔").
			WithDatabase("datab'ase").
			WithShardDurationString("1h").
			WithDurationString("3h").
			WithReplicationFactor(36),
		"CREATE RETENTION POLICY \"PolicyName↔\" ON \"datab'ase\" DURATION 3h REPLICATION 36 SHARD DURATION 1h",
		false,
	},
	{
		"Simple Retention policy invalid duration",
		NewCreateRetentionPolicyBuilder().
			WithPolicyName("PolicyName↔").
			WithDatabase("datab'ase").
			WithShardDurationString("notValid").
			WithDurationString("notValid").
			WithReplicationFactor(36),
		"CREATE RETENTION POLICY \"PolicyName↔\" ON \"datab'ase\" DURATION 0s REPLICATION 36",
		false,
	},
	{
		"Retention policy with future and past limit",
		NewCreateRetentionPolicyBuilder().WithPolicyName("policy1").
			WithDatabase("testdb").
			WithDurationString("1h").
			WithReplicationFactor(2).
			WithShardDurationString("1s").
			WithFutureLimitString("12h").
			WithPastLimitString("3s"),
		"CREATE RETENTION POLICY policy1 ON testdb DURATION 1h REPLICATION 2 SHARD DURATION 1s FUTURE LIMIT 12h PAST LIMIT 3s",
		false,
	},
}

func TestCreateRetentionPolicyBuilder(t *testing.T) {
	for i, sample := range testCreateRetentionPolicy {
		s, err := sample.b.Build()

		fmt.Print("Test ", i, ": ", sample.d)

		if sample.e {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		assert.Equal(t, sample.s, s)

		fmt.Println("   [OK]")
	}
}
