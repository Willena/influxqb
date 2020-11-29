package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testExplainBuilder = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{
		"Explain select",
		NewExplainBuilder().
			WithSelectBuilder(NewSelectBuilder().Select("F1").From("Toto").Where(Eq(Field{Name: "totoField"}, 32))),
		"EXPLAIN SELECT F1 FROM Toto WHERE (totoField = 32)",
		false,
	},
	{
		"Explain analyse select",
		NewExplainBuilder().Analyze().
			WithSelectBuilder(NewSelectBuilder().Select("F1").From("Toto").Where(Eq(Field{Name: "totoField"}, 32))),
		"EXPLAIN ANALYZE SELECT F1 FROM Toto WHERE (totoField = 32)",
		false,
	},
}

func TestExplainBuilder(t *testing.T) {
	for i, sample := range testExplainBuilder {
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
