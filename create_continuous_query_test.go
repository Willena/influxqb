package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCreateContinuousQueryBuilder = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{
		"Simple Continuous Query ",
		NewCreateContinuousQueryBuilder().WithName("ContinuousQuery").WithDatabase("db").
			WithResamplingIntervalFromString("12h").WithTimeoutString("1h"),
		"CREATE CONTINUOUS QUERY ContinuousQuery ON db RESAMPLE EVERY 12h FOR 1h BEGIN  END",
		false,
	},
	{
		"Simple Continuous Query with select builder ",
		NewCreateContinuousQueryBuilder().WithName("ContinuousQuery").WithDatabase("db").
			WithResamplingIntervalFromString("12h").WithTimeoutString("1h").
			WithSelectBuilder(NewSelectBuilder().Select("F1").From("Toto").Where(Eq(Field{Name: "totoField"}, 32))),
		"CREATE CONTINUOUS QUERY ContinuousQuery ON db RESAMPLE EVERY 12h FOR 1h BEGIN SELECT F1 FROM Toto WHERE (totoField = 32) END",
		false,
	},
}

func TestCreateContinuousQueryBuilder(t *testing.T) {
	for i, sample := range testCreateContinuousQueryBuilder {
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
