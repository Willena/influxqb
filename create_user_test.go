package influxqb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCreateUser = []struct {
	d string
	b BuilderIf
	s string
	e bool
}{
	{"Create User non Admin",
		NewCreateUserBuilder().WithUsername("UserName").WithPassword("password!&é!$^ù"),
		"CREATE USER UserName WITH PASSWORD 'password!&é!$^ù'",
		false,
	},
	{"Create User Admin",
		NewCreateUserBuilder().WithUsername("UserName").WithPassword("password!&é!$^ù").Admin(),
		"CREATE USER UserName WITH PASSWORD 'password!&é!$^ù' WITH ALL PRIVILEGES",
		false,
	},
}

func TestCreateUserBuilder(t *testing.T) {
	for i, sample := range testCreateUser {
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
