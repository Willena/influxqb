package influxqb

import (
	"fmt"
	"github.com/influxdata/influxql"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

/*


 */

var testDropSeriesBuilder = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"Drop Series from ",
		NewDropSeries().From(NewMeasurement().Name("M")),
		"DROP SERIES FROM M",
		false,
	},
	{"Drop Series from ",
		NewDropSeries().From(NewMeasurement().Name("M").WithDatabase("db").WithPolicy("policy1")),
		"DROP SERIES FROM db.policy1.M",
		false,
	},
	{"Drop Series from Where",
		NewDropSeries().
			From(NewMeasurement().Regex(regexp.MustCompile(".*"))).
			Where(Eq(&Field{Name: "cpu"}, "cpu2")),
		"DROP SERIES FROM /.*/ WHERE (cpu = 'cpu2')",
		false,
	},
	{"Drop series where ",
		NewDropSeries().
			Where(Eq(&Field{Name: "cpu"}, "cpu2")),
		"DROP SERIES WHERE (cpu = 'cpu2')",
		false,
	},
	{"Drop Series MathExpr",
		NewDropSeries().
			Where(&Math{Expr: []interface{}{
				&Field{Name: "cpu"}, influxql.EQ, "cpu2"}},
			),
		"DROP SERIES WHERE (cpu = 'cpu2')",
		false,
	},
}

func TestDropSeriesBuilder(t *testing.T) {
	for i, sample := range testDropSeriesBuilder {
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
