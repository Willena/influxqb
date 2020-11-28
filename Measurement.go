package influxqb

import (
	"github.com/influxdata/influxql"
	"regexp"
)

func NewMeasurement() *Measurement {
	return &Measurement{m: &influxql.Measurement{}}
}

type Measurement struct {
	m *influxql.Measurement
}

func (m *Measurement) Name(str string) *Measurement {
	m.m.Name = str
	return m
}

func (m *Measurement) Regex(regex *regexp.Regexp) *Measurement {
	m.m.Regex = &influxql.RegexLiteral{Val: regex}
	return m
}

func (m *Measurement) WithDatabase(str string) *Measurement {
	m.m.Database = str
	return m
}
func (m *Measurement) WithPolicy(str string) *Measurement {
	m.m.RetentionPolicy = str
	return m
}
