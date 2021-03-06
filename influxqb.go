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

func NewKillQueryBuilder() *KillQueryBuilder {
	return &KillQueryBuilder{dss: &influxql.KillQueryStatement{}}
}

func NewShowContinuousQueries() *ShowContinuousQueries {
	return &ShowContinuousQueries{}
}

func NewShowDatabases() *ShowDatabases {
	return &ShowDatabases{}
}

func NewShowUsers() *ShowUsers {
	return &ShowUsers{}
}

func NewShowQueries() *ShowQueries {
	return &ShowQueries{}
}

func NewShowShardGroup() *ShowShardGroups {
	return &ShowShardGroups{}
}

func NewShowSubscriptions() *ShowSubscriptions {
	return &ShowSubscriptions{}
}

func NewShowShards() *ShowShards {
	return &ShowShards{}
}

func NewShowDiagnostics() *ShowDiagnostics {
	return &ShowDiagnostics{}
}

func NewShowFieldKeyCardinality() *ShowFieldKeyCardinality {
	return &ShowFieldKeyCardinality{q: &influxql.ShowFieldKeyCardinalityStatement{}}
}

func NewShowFieldKeys() *ShowFieldKeys {
	return &ShowFieldKeys{q: &influxql.ShowFieldKeysStatement{}}
}

func NewShowGrants() *ShowGrants {
	return &ShowGrants{q: &influxql.ShowGrantsForUserStatement{}}
}

func NewShowMeasurementCardinality() *ShowMeasurementCardinality {
	return &ShowMeasurementCardinality{q: &influxql.ShowMeasurementCardinalityStatement{}}
}

func NewShowMeasurements() *ShowMeasurements {
	return &ShowMeasurements{q: &influxql.ShowMeasurementsStatement{}}
}

func NewShowRetentionPolicies() *ShowRetentionPolicies {
	return &ShowRetentionPolicies{}
}

func NewShowSeries() *ShowSeries {
	return &ShowSeries{q: &influxql.ShowSeriesStatement{}}
}

func NewShowSeriesCadinality() *ShowSeriesCadinality {
	return &ShowSeriesCadinality{q: &influxql.ShowSeriesCardinalityStatement{}}
}

func NewShowStats() *ShowStats {
	return &ShowStats{q: &influxql.ShowStatsStatement{}}
}

func NewShowTagKeyCardinality() *ShowTagKeyCardinality {
	return &ShowTagKeyCardinality{q: &influxql.ShowTagKeyCardinalityStatement{}}
}

func NewShowTagKeys() *ShowTagKeys {
	return &ShowTagKeys{q: &influxql.ShowTagKeysStatement{}}
}

func NewShowTagValues() *ShowTagValues {
	return &ShowTagValues{q: &influxql.ShowTagValuesStatement{}}
}

func NewShowTagValuesCardinality() *ShowTagValuesCardinality {
	return &ShowTagValuesCardinality{q: &influxql.ShowTagValuesCardinalityStatement{}}
}
