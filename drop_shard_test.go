package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*


 */

var testDropShardBuilder = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"DROP Shard",
		NewDropShard().WithShard(1),
		"DROP SHARD 1",
		false,
	},
}

func TestDropShardBuilder(t *testing.T) {
	for i, sample := range testDropShardBuilder {
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
