package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*


 */

var testDropMeasurementBuilder = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"DROP continuous query named",
		NewDropMeasurement().WithMeasurement("Measurement"),
		"DROP MEASUREMENT \"Measurement\"",
		false,
	},
}

func TestDropMeasurementBuilder(t *testing.T) {
	for i, sample := range testDropMeasurementBuilder {
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
