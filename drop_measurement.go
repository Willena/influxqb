package influxqb

import "github.com/influxdata/influxql"

type DropMeasurementBuilder struct {
	dcq *influxql.DropMeasurementStatement
}

func (b *DropMeasurementBuilder) WithMeasurement(Measurement string) *DropMeasurementBuilder {
	b.dcq.Name = Measurement
	return b
}

func (b *DropMeasurementBuilder) Build() (string, error) {
	return b.dcq.String(), nil
}
