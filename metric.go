package godruid

type MetricType string
type OrderType string

var (
	// Metric
	NumericMetricType   MetricType = "numeric"
	DimensionMetricType MetricType = "dimension"
	InvertedMetricType  MetricType = "inverted"
	// Ordering
	LexicographicOrder OrderType = "lexicographic"
	AlphaNumericOrder  OrderType = "alphanumeric"
	NumericOrder       OrderType = "numeric"
	StrlenOrder        OrderType = "strlen"
)

type Metric struct {
	Type         MetricType `json:"type"`
	MetricName   string     `json:"metric,omitempty"`
	Ordering     OrderType  `json:"ordering,omitempty"`
	PreviousStop string     `json:"previousStop,omitempty"`
	// This is another annoying one, probably a custom Marshal
	//Metric       *Metric    `json:"metric,omitempty"`
}

func NewDimensionOrderedMetric(ot OrderType) *Metric {
	return &Metric{
		Type:     DimensionMetricType,
		Ordering: ot,
	}
}

func NewNumericOrderedMetric(metricName string) *Metric {
	return &Metric{
		Type:       NumericMetricType,
		MetricName: metricName,
	}
}
