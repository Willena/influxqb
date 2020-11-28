package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*


 */

var testDropDatabaseBuilder = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"DROP continuous query named",
		NewDropDatabase().WithDatabase("database"),
		"DROP DATABASE \"database\"",
		false,
	},
}

func TestDropDatabaseBuilder(t *testing.T) {
	for i, sample := range testDropDatabaseBuilder {
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
