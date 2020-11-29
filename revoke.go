package influxqb

import "github.com/influxdata/influxql"

type RevokeBuilder struct {
	isAdmin   bool
	db        string
	username  string
	privilege influxql.Privilege
}

func (b *RevokeBuilder) Admin() *RevokeBuilder {
	b.isAdmin = true
	return b
}

func (b *RevokeBuilder) FromUsername(username string) *RevokeBuilder {
	b.username = username
	return b
}

func (b *RevokeBuilder) OnDatabase(database string) *RevokeBuilder {
	b.db = database
	return b
}

func (b *RevokeBuilder) WithPrivilege(privilege influxql.Privilege) *RevokeBuilder {
	b.privilege = privilege
	return b
}

func (b *RevokeBuilder) Build() (string, error) {

	if b.isAdmin {
		stm := &influxql.RevokeAdminStatement{User: b.username}
		return stm.String(), nil
	} else {
		stm := &influxql.RevokeStatement{
			Privilege: b.privilege,
			On:        b.db,
			User:      b.username,
		}
		return stm.String(), nil
	}

}
