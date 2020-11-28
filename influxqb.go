package influxqb

import "github.com/influxdata/influxql"

type BuilderIf interface {
	Build() (string, error)
}

func NewSelectBuilder() *SelectBuilder {
	return &SelectBuilder{selectStatement: &influxql.SelectStatement{}}
}

func NewAlterRetentionPolicyBuilder() *AlterRetentionPolicyBuilder {
	return &AlterRetentionPolicyBuilder{alterStm: &influxql.AlterRetentionPolicyStatement{}}
}

func NewCreateContinuousQueryBuilder() *CreateContinuousQueryBuilder {
	return &CreateContinuousQueryBuilder{continuousQuery: &influxql.CreateContinuousQueryStatement{}, selectBuilder: NewSelectBuilder()}
}

func NewCreateRetentionPolicyBuilder() *CreateRetentionPolicyBuilder {
	return &CreateRetentionPolicyBuilder{ret: &influxql.CreateRetentionPolicyStatement{}}
}

func NewCreateDatabaseBuilder() *CreateDatabaseBuilder {
	return &CreateDatabaseBuilder{dbStatement: &influxql.CreateDatabaseStatement{}}
}

func NewCreateSubscriptionBuilder() *CreateSubscriptionBuilder {
	return &CreateSubscriptionBuilder{subStm: &influxql.CreateSubscriptionStatement{}}
}

func NewCreateUserBuilder() *CreateUserBuilder {
	return &CreateUserBuilder{cu: &influxql.CreateUserStatement{}}
}

func NewDeleteBuilder() *DeleteBuilder {
	return &DeleteBuilder{del: &influxql.DeleteStatement{}}
}

func NewDropContinuousQuery() *DropContinuousQueryBuilder {
	return &DropContinuousQueryBuilder{dcq: &influxql.DropContinuousQueryStatement{}}
}

func NewDropDatabase() *DropDatabaseBuilder {
	return &DropDatabaseBuilder{dcq: &influxql.DropDatabaseStatement{}}
}

func NewDropMeasurement() *DropMeasurementBuilder {
	return &DropMeasurementBuilder{dcq: &influxql.DropMeasurementStatement{}}
}
