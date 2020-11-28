package influxqb

import (
	"bytes"
	"github.com/influxdata/influxql"
)

type CreateUserBuilder struct {
	cu *influxql.CreateUserStatement
}

func (b *CreateUserBuilder) WithUsername(str string) *CreateUserBuilder {
	b.cu.Name = str
	return b
}

func (b *CreateUserBuilder) WithPassword(str string) *CreateUserBuilder {
	b.cu.Password = str
	return b
}

func (b *CreateUserBuilder) Admin() *CreateUserBuilder {
	b.cu.Admin = true
	return b
}

func (b *CreateUserBuilder) Build() (string, error) {

	// The InfluxQL parser automaticaly replace the password by [REDACTED] and hide the password
	// We are forced to rewirte the String() function by ourself.

	var buf bytes.Buffer
	_, _ = buf.WriteString("CREATE USER ")
	_, _ = buf.WriteString(influxql.QuoteIdent(b.cu.Name))
	_, _ = buf.WriteString(" WITH PASSWORD ")
	_, _ = buf.WriteString(influxql.QuoteString(b.cu.Password))
	if b.cu.Admin {
		_, _ = buf.WriteString(" WITH ALL PRIVILEGES")
	}
	return buf.String(), nil
}
