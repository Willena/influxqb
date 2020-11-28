package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*


 */

var testDropRetentionPolicyBuilder = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"DROP continuous query named",
		NewDropRetentionPolicy().WithRetentionPolicy("RetentionPolicy").WithDatabase("database"),
		"DROP RETENTION POLICY RetentionPolicy ON \"database\"",
		false,
	},
}

func TestDropRetentionPolicyBuilder(t *testing.T) {
	for i, sample := range testDropRetentionPolicyBuilder {
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
