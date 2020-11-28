package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testSubscription = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"Create Subscription ALL single host",
		NewCreateSubscriptionBuilder().WithName("Name").
			WithDatabase("db").WithRetentionPolicy("policy").WithMode(ALL).WithDestination("host1"),
		"CREATE SUBSCRIPTION \"Name\" ON db.\"policy\" DESTINATIONS ALL 'host1'",
		false,
	},
	{"Create Subscription ANY single host",
		NewCreateSubscriptionBuilder().WithName("Name").
			WithDatabase("db").WithRetentionPolicy("policy").WithMode(ANY).WithDestination("host1"),
		"CREATE SUBSCRIPTION \"Name\" ON db.\"policy\" DESTINATIONS ANY 'host1'",
		false,
	},
	{"Create Subscription ANY multi host",
		NewCreateSubscriptionBuilder().WithName("Name").
			WithDatabase("db").WithRetentionPolicy("policy").WithMode(ANY).WithDestinations("host1", "host2"),
		"CREATE SUBSCRIPTION \"Name\" ON db.\"policy\" DESTINATIONS ANY 'host1', 'host2'",
		false,
	},
}

func TestCreateSubscriptionBuilder(t *testing.T) {
	for i, sample := range testSubscription {
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
