package influxqb

import (
	"fmt"
	"github.com/influxdata/influxql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/*


 */

var testDeleteBuilder = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"Delete all from mesurement",
		NewDeleteBuilder().From("cpu"),
		"DELETE FROM cpu",
		false,
	},
	{"Delete from mesurement where ",
		NewDeleteBuilder().
			From("cpu").
			Where(LessThan(&Field{Name: "time"}, time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC))),
		"DELETE FROM cpu WHERE (time < '2000-01-01T00:00:00Z')",
		false,
	},
	{"Delete all WHERE time < 2000",
		NewDeleteBuilder().
			Where(LessThan(&Field{Name: "time"}, time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC))),
		"DELETE FROM /.*/ WHERE (time < '2000-01-01T00:00:00Z')",
		false,
	},
	{"Delete all WHERE time < 2000 with MathExpr",
		NewDeleteBuilder().
			Where(&Math{Expr: []interface{}{
				&Field{Name: "time"}, influxql.LT, time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC)}},
			),
		"DELETE FROM /.*/ WHERE (time < '2000-01-01T00:00:00Z')",
		false,
	},
}

func TestDeleteBuilder(t *testing.T) {
	for i, sample := range testDeleteBuilder {
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
