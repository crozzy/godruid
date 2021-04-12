package godruid

import (
	"encoding/json"
)

type Aggregator struct {
	AggType    string      `json:"type"`
	Name       string      `json:"name,omitempty"`
	FieldName  string      `json:"fieldName,omitempty"`
	Filter     *Filter     `json:"filter,omitempty"`
	Aggregator *Aggregator `json:"aggregator,omitempty"`
}

func (a *Aggregator) toJSON() ([]byte, error) {
	return json.Marshal(a)
}

func NewHyperUniqueAggregator(metricName, outputName string) *Aggregator {
	return &Aggregator{AggType: "hyperUnique", Name: outputName, FieldName: metricName}
}

func NewDoubleSumAggregator(metricName, outputName string) *Aggregator {
	return &Aggregator{AggType: "doubleSum", Name: outputName, FieldName: metricName}
}

func NewCountAggregator(outputName string) *Aggregator {
	return &Aggregator{AggType: "count", Name: outputName}
}

func NewFilteredAggregator(name string, filter *Filter, agg *Aggregator) *Aggregator {
	return &Aggregator{AggType: "filtered", Name: name, Filter: filter, Aggregator: agg}
}

func NewThetaSketchAggregator(metricName, outputName string) *Aggregator {
	return &Aggregator{AggType: "thetaSketch", Name: outputName, FieldName: metricName}
}
