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
