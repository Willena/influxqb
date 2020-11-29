package influxqb

import (
	"fmt"
	"github.com/influxdata/influxql"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*


 */

var testRevokeBuilder = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"Revoke Read on all for user",
		NewRevokeBuilder().FromUsername("Usr").WithPrivilege(influxql.ReadPrivilege).OnDatabase("db"),
		"REVOKE READ ON db FROM Usr",
		false,
	},
	{"Revoke write on db for user",
		NewRevokeBuilder().FromUsername("Usr").WithPrivilege(influxql.WritePrivilege).OnDatabase("db"),
		"REVOKE WRITE ON db FROM Usr",
		false,
	},
	{"Revoke All on db for user ",
		NewRevokeBuilder().FromUsername("Usr").WithPrivilege(influxql.AllPrivileges).OnDatabase("db"),
		"REVOKE ALL PRIVILEGES ON db FROM Usr",
		false,
	},
	{"Revoke no privilege on db ",
		NewRevokeBuilder().FromUsername("Usr").WithPrivilege(influxql.NoPrivileges).OnDatabase("db"),
		"REVOKE NO PRIVILEGES ON db FROM Usr",
		false,
	},
	{"Revoke Admin on db",
		NewRevokeBuilder().Admin().FromUsername("Usr"),
		"REVOKE ALL PRIVILEGES FROM Usr",
		false,
	},
}

func TestRevokeBuilder(t *testing.T) {
	for i, sample := range testRevokeBuilder {
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
