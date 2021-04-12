package godruid

import (
	"encoding/json"
)

type query interface {
	ToJSON() ([]byte, error)
}

type BaseQuery struct {
	QueryType       string           `json:"queryType"`
	DataSource      string           `json:"dataSource"`
	Granularity     Granularity      `json:"granularity"`
	Intervals       []string         `json:"intervals"`
	Aggregations    []*Aggregator    `json:"aggregations"`
	Dimensions      []string         `json:"dimensions"`
	Filter          *Filter          `json:"filter"`
	Context         *QueryContext    `json:"context,omitempty"`
	PostAggregation *PostAggregation `json:"postAggregation,omitempty"`
}

func (q *BaseQuery) AddFilter(filter *Filter) {
	q.Filter = filter
}

func (q *BaseQuery) AddContext(context *QueryContext) {
	q.Context = context
}

func (q *BaseQuery) AddPostAggregation(pp *PostAggregation) {
	q.PostAggregation = pp
}

func (q *BaseQuery) ToJSON() ([]byte, error) {
	return json.MarshalIndent(q, "", "  ")
}

func (q *BaseQuery) GetContext() *QueryContext {
	return q.Context
}

// NewGroupByQuery
func NewGroupByQuery(
	dataSource string,
	granularity Granularity,
	intervals []string,
	aggregations []*Aggregator,
	dimensions []string,
) *GroupByQuery {
	q := &GroupByQuery{}
	q.DataSource = dataSource
	q.Granularity = granularity
	q.Intervals = intervals
	q.Aggregations = aggregations
	q.Dimensions = dimensions
	q.Filter = &Filter{}
	q.QueryType = "groupBy"
	return q
}

type GroupByQuery struct {
	BaseQuery
	LimitSpec *LimitSpec `json:"limitSpec,omitempty"`
}

func (q *GroupByQuery) AddLimit(limitSpec *LimitSpec) {
	q.LimitSpec = limitSpec
}

// NewTopNQuery
func NewTopNQuery(
	dataSource string,
	intervals []string,
	dimension *Dimension,
	metric *Metric,
	aggregations []*Aggregator,
	threshold int64,
) *TopNQuery {
	q := &TopNQuery{}
	q.DataSource = dataSource
	q.Intervals = intervals
	q.Aggregations = aggregations
	q.Dimension = dimension
	q.Metric = metric
	q.Filter = &Filter{}
	q.QueryType = "topN"
	q.Threshold = threshold
	q.Granularity = NewAllGranularity()
	return q
}

type TopNQuery struct {
	BaseQuery
	Metric    *Metric    `json:"metric"`
	Dimension *Dimension `json:"dimension"`
	Threshold int64      `json:"threshold"`
}

// NewTimeseriesQuery
func NewTimeseriesQuery(
	dataSource string,
	intervals []string,
	aggregations []*Aggregator,
	granularity Granularity,
	decending bool,
) *TimeseriesQuery {
	q := &TimeseriesQuery{
		DataSource:   dataSource,
		Intervals:    intervals,
		Aggregations: aggregations,
		Filter:       &Filter{},
		QueryType:    "timeseries",
		Granularity:  granularity,
		Decending:    decending,
	}
	return q
}

type TimeseriesQuery struct {
	QueryType        string             `json:"queryType"`
	DataSource       string             `json:"dataSource"`
	Decending        bool               `json:"decending"`
	Intervals        []string           `json:"intervals"`
	Granularity      Granularity        `json:"granularity"`
	Filter           *Filter            `json:"filter"`
	Aggregations     []*Aggregator      `json:"aggregations"`
	PostAggregations []*PostAggregation `json:"postAggregations,omitempty"`
	Context          *QueryContext      `json:"context,omitempty"`
}

func (tsq *TimeseriesQuery) AddFilter(filter *Filter) {
	tsq.Filter = filter
}

func (tsq *TimeseriesQuery) AddContext(context *QueryContext) {
	tsq.Context = context
}

func (tsq *TimeseriesQuery) AddPostAggregation(pp *PostAggregation) {
	tsq.PostAggregations = append(tsq.PostAggregations, pp)
}

func (tsq *TimeseriesQuery) ToJSON() ([]byte, error) {
	return json.MarshalIndent(tsq, "", "  ")
}
