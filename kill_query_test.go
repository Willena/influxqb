package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*


 */

var testKillQueryBuilder = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"Kill query",
		NewKillQueryBuilder().WithQueryId(1),
		"KILL QUERY 1",
		false,
	},
	{"kill query on host ",
		NewKillQueryBuilder().WithQueryId(1).OnHost("hostname"),
		"KILL QUERY 1 ON hostname",
		false,
	},
}

func TestKillQueryBuilder(t *testing.T) {
	for i, sample := range testKillQueryBuilder {
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
