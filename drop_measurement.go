package influxqb

import "github.com/influxdata/influxql"

type DropMeasurementBuilder struct {
	dms *influxql.DropMeasurementStatement
}

func (b *DropMeasurementBuilder) WithMeasurement(Measurement string) *DropMeasurementBuilder {
	b.dms.Name = Measurement
	return b
}

func (b *DropMeasurementBuilder) Build() (string, error) {
	return b.dms.String(), nil
}
