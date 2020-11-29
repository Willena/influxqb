package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*


 */

var testDropUserBuilder = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"DROP User",
		NewDropUser().WithUsername("Usr"),
		"DROP USER Usr",
		false,
	},
}

func TestDropUserBuilder(t *testing.T) {
	for i, sample := range testDropUserBuilder {
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
