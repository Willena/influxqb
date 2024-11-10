package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*


 */

var testDropSubscriptionBuilder = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"DROP continuous query named",
		NewDropSubscription().
			WithName("Subscription").
			WithDatabase("database").
			WithRetentionPolicy("policy"),
		"DROP SUBSCRIPTION \"Subscription\" ON \"database\".\"policy\"",
		false,
	},
}

func TestDropSubscriptionBuilder(t *testing.T) {
	for i, sample := range testDropSubscriptionBuilder {
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
