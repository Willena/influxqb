package influxqb

import "github.com/influxdata/influxql"

type GrantBuilder struct {
	isAdmin   bool
	db        string
	username  string
	privilege influxql.Privilege
}

func (b *GrantBuilder) Admin() *GrantBuilder {
	b.isAdmin = true
	return b
}

func (b *GrantBuilder) WithUsername(username string) *GrantBuilder {
	b.username = username
	return b
}

func (b *GrantBuilder) WithDatabase(database string) *GrantBuilder {
	b.db = database
	return b
}

func (b *GrantBuilder) WithPrivilege(privilege influxql.Privilege) *GrantBuilder {
	b.privilege = privilege
	return b
}

func (b *GrantBuilder) Build() (string, error) {

	if b.isAdmin {
		stm := &influxql.GrantAdminStatement{User: b.username}
		return stm.String(), nil
	} else {
		stm := &influxql.GrantStatement{
			Privilege: b.privilege,
			On:        b.db,
			User:      b.username,
		}
		return stm.String(), nil
	}

}
