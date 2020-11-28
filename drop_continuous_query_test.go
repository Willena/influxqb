package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*


 */

var testDropContinuousQueryBuilder = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"DROP continuous query named",
		NewDropContinuousQuery().WithName("QueryName").WithDatabase("database"),
		"DROP CONTINUOUS QUERY QueryName ON \"database\"",
		false,
	},
}

func TestDropContinuousQueryBuilder(t *testing.T) {
	for i, sample := range testDropContinuousQueryBuilder {
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
