package influxqb

import (
	"fmt"
	"github.com/influxdata/influxql"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*


 */

var testGrantBuilder = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"Grant Read on all for user",
		NewGrantBuilder().
			WithUsername("Usr").
			WithPrivilege(influxql.ReadPrivilege).
			WithDatabase("db"),
		"GRANT READ ON db TO Usr",
		false,
	},
	{"grant write on db for user",
		NewGrantBuilder().
			WithUsername("Usr").
			WithPrivilege(influxql.WritePrivilege).
			WithDatabase("db"),
		"GRANT WRITE ON db TO Usr",
		false,
	},
	{"Grant All on db for user ",
		NewGrantBuilder().
			WithUsername("Usr").
			WithPrivilege(influxql.AllPrivileges).
			WithDatabase("db"),
		"GRANT ALL PRIVILEGES ON db TO Usr",
		false,
	},
	{"Grant no privilege on db ",
		NewGrantBuilder().
			WithUsername("Usr").
			WithPrivilege(influxql.NoPrivileges).
			WithDatabase("db"),
		"GRANT NO PRIVILEGES ON db TO Usr",
		false,
	},
	{"Grant Admin on db",
		NewGrantBuilder().Admin().WithUsername("Usr"),
		"GRANT ALL PRIVILEGES TO Usr",
		false,
	},
}

func TestGrantBuilder(t *testing.T) {
	for i, sample := range testGrantBuilder {
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
