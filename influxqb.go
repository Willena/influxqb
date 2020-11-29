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
	return &DropMeasurementBuilder{dms: &influxql.DropMeasurementStatement{}}
}

func NewDropRetentionPolicy() *DropRetentionPolicyBuilder {
	return &DropRetentionPolicyBuilder{drp: &influxql.DropRetentionPolicyStatement{}}
}

func NewDropSeries() *DropSeriesBuilder {
	return &DropSeriesBuilder{dss: &influxql.DropSeriesStatement{}}
}

func NewDropShard() *DropShardBuilder {
	return &DropShardBuilder{dss: &influxql.DropShardStatement{}}
}

func NewDropSubscription() *DropSubscriptionBuilder {
	return &DropSubscriptionBuilder{dss: &influxql.DropSubscriptionStatement{}}
}

func NewDropUser() *DropUserBuilder {
	return &DropUserBuilder{dss: &influxql.DropUserStatement{}}
}

func NewExplainBuilder() *ExplainBuilder {
	return &ExplainBuilder{explainStatement: &influxql.ExplainStatement{}, selectBuilder: &SelectBuilder{}}
}

func NewGrantBuilder() *GrantBuilder {
	return &GrantBuilder{}
}

func NewRevokeBuilder() *RevokeBuilder {
	return &RevokeBuilder{}
}
