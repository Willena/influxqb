package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testAlterRetentionPolicy = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{
		"Simple Retention policy",
		NewAlterRetentionPolicyBuilder().WithPolicyName("PolicyName").WithDatabase("database"),
		"ALTER RETENTION POLICY PolicyName ON \"database\"",
		false,
	},
	{
		"Simple Retention policy set default",
		NewAlterRetentionPolicyBuilder().
			WithPolicyName("PolicyName").
			WithDatabase("database").
			SetAsDefault(),
		"ALTER RETENTION POLICY PolicyName ON \"database\" DEFAULT",
		false,
	},
	{
		"Simple Retention policy set duration",
		NewAlterRetentionPolicyBuilder().
			WithPolicyName("PolicyName").
			WithDatabase("database").
			WithDurationString("1h"),
		"ALTER RETENTION POLICY PolicyName ON \"database\" DURATION 1h",
		false,
	},
	{
		"Simple Retention policy set  shard duration",
		NewAlterRetentionPolicyBuilder().
			WithPolicyName("PolicyName").
			WithDatabase("database").
			WithShardDurationString("1h"),
		"ALTER RETENTION POLICY PolicyName ON \"database\" SHARD DURATION 1h",
		false,
	},
	{
		"Simple Retention policy set  shard duration and duration",
		NewAlterRetentionPolicyBuilder().
			WithPolicyName("PolicyName").
			WithDatabase("database").
			WithShardDurationString("1h").
			WithDurationString("3h"),
		"ALTER RETENTION POLICY PolicyName ON \"database\" DURATION 3h SHARD DURATION 1h",
		false,
	},
	{
		"Simple Retention policy set  shard duration and duration plus replication factor",
		NewAlterRetentionPolicyBuilder().
			WithPolicyName("PolicyName").
			WithDatabase("database").
			WithShardDurationString("1h").
			WithDurationString("3h").
			WithReplicationFactor(36),
		"ALTER RETENTION POLICY PolicyName ON \"database\" DURATION 3h REPLICATION 36 SHARD DURATION 1h",
		false,
	},
	{
		"Simple Retention policy unicode names and quotes",
		NewAlterRetentionPolicyBuilder().
			WithPolicyName("PolicyName↔").
			WithDatabase("datab'ase").
			WithShardDurationString("1h").
			WithDurationString("3h").
			WithReplicationFactor(36),
		"ALTER RETENTION POLICY \"PolicyName↔\" ON \"datab'ase\" DURATION 3h REPLICATION 36 SHARD DURATION 1h",
		false,
	},
	{
		"Simple Retention policy invalid duration",
		NewAlterRetentionPolicyBuilder().
			WithPolicyName("PolicyName↔").
			WithDatabase("datab'ase").
			WithShardDurationString("notValid").
			WithDurationString("notValid").
			WithReplicationFactor(36),
		"ALTER RETENTION POLICY \"PolicyName↔\" ON \"datab'ase\" DURATION 0s REPLICATION 36 SHARD DURATION 0s",
		false,
	},
	{
		"Alter with future and past limits",
		NewAlterRetentionPolicyBuilder().
			WithPolicyName("default").
			WithDatabase("testdb").
			WithDurationString("0s").
			WithReplicationFactor(1).
			WithFutureLimitString("55s").
			WithPastLimitString("52m"),
		"ALTER RETENTION POLICY \"default\" ON testdb DURATION 0s REPLICATION 1 FUTURE LIMIT 55s PAST LIMIT 52m",
		false,
	},
}

func TestAlterRetentionPolicyBuilder(t *testing.T) {
	for i, sample := range testAlterRetentionPolicy {
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
