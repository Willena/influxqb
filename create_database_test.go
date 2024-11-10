package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCreateDatabase = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"Create Simple database no policy",
		NewCreateDatabaseBuilder().WithName("DatabaseName"),
		"CREATE DATABASE DatabaseName",
		false,
	},
	{"Create Simple database With policy",
		NewCreateDatabaseBuilder().WithName("DatabaseName").WithRetentionPolicy(
			NewCreateRetentionPolicyBuilder().
				WithReplicationFactor(3).
				WithDurationString("1h").
				WithShardDurationString("2h"),
		),
		"CREATE DATABASE DatabaseName WITH DURATION 1h0m0s REPLICATION 3 SHARD DURATION 2h0m0s",
		false,
	},
	{
		"Create database with future and past policy",
		NewCreateDatabaseBuilder().
			WithName("testdb").
			WithRetentionPolicy(
				NewCreateRetentionPolicyBuilder().
					WithPolicyName("test_name").
					WithDurationString("24h").
					WithShardDurationString("10m").
					WithReplicationFactor(2).
					WithPastLimitString("67ms").
					WithFutureLimitString("1h")),
		"CREATE DATABASE testdb WITH DURATION 24h0m0s REPLICATION 2 SHARD DURATION 10m0s FUTURE LIMIT 1h PAST LIMIT 67ms NAME test_name",
		false,
	},
}

func TestCreateDatabaseBuilder(t *testing.T) {
	for i, sample := range testCreateDatabase {
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
