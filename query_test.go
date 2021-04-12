package godruid

import (
	"fmt"
	"testing"
	"time"
)

func TestGroupBy(t *testing.T) {
	origin, _ := time.Parse(time.RFC3339, "2020-01-20T00:00:00.00Z")
	granularity := NewPeriodGranularity(P1D, "Europe/Prague", origin)
	query := NewGroupByQuery("test", granularity, []string{"something"}, []*Aggregator{NewCountAggregator("thing")}, []string{"bloop"})
	res, err := query.ToJSON()
	if err != nil {
		t.Error("got an error", err)
	}
	t.Log(string(res))
}

func TestGroupByWithFilters(t *testing.T) {
	origin, _ := time.Parse(time.RFC3339, "2020-01-20T00:00:00.00Z")
	granularity := NewPeriodGranularity(P1D, "Europe/Prague", origin)
	query := NewGroupByQuery("test", granularity, []string{"2020-01-20T00:00:00.000/2020-01-21T00:00:00.000"}, []*Aggregator{NewCountAggregator("thing")}, []string{"application_id"})

	// Create filters
	filter := NewBaseFilter(OrExp)
	appIDFilter := NewDimensionFilter("application_id", "1347123")
	appIDFilter2 := NewDimensionFilter("application_id", "1010123")
	filter.AddChildren(appIDFilter, appIDFilter2)
	query.AddFilter(filter)
	res, err := query.ToJSON()
	if err != nil {
		t.Error("got an error", err)
	}
	t.Log(string(res))
}

func TestGroupByWithFilteredAggregator(t *testing.T) {
	origin, _ := time.Parse(time.RFC3339, "2020-01-20T00:00:00.00Z")
	granularity := NewPeriodGranularity(P1D, "Europe/Prague", origin)
	aggFilter := NewInFilter("event_name", []string{"Purchase"})
	revAggregator := NewDoubleSumAggregator("dim_sum", "Revenue")
	filteredAgg := NewFilteredAggregator("Revenue", aggFilter, revAggregator)
	query := NewGroupByQuery("test", granularity, []string{"2020-01-20T00:00:00.000/2020-01-21T00:00:00.000"}, []*Aggregator{NewCountAggregator("thing"), filteredAgg}, []string{"application_id"})

	// Create filters
	filter := NewBaseFilter(OrExp)
	appIDFilter := NewDimensionFilter("application_id", "1347123")
	appIDFilter2 := NewDimensionFilter("application_id", "1010123")
	filter.AddChildren(appIDFilter, appIDFilter2)
	query.AddFilter(filter)
	res, err := query.ToJSON()
	if err != nil {
		t.Error("got an error", err)
	}
	t.Log(string(res))
}

func TestGroupByNoGranularity(t *testing.T) {
	granularity := NewAllGranularity()
	query := NewGroupByQuery("test", granularity, []string{"something"}, []*Aggregator{NewCountAggregator("thing")}, []string{"bloop"})
	res, err := query.ToJSON()
	if err != nil {
		t.Error("got an error", err)
	}
	t.Log(string(res))
}

func TestTopN(t *testing.T) {
	query := NewTopNQuery("test", []string{"something"}, NewDefaultDimension("application_id", "application_id", "string"), NewNumericOrderedMetric("count"), []*Aggregator{NewCountAggregator("thing")}, 10)
	res, err := query.ToJSON()
	if err != nil {
		t.Error("got an error", err)
	}
	t.Log(string(res))
}

func TestCustomDimensionsTopN(t *testing.T) {
	query := NewTopNQuery(
		"test",
		[]string{"something"},
		NewPartialExtractionDimension("custom_properties", "this", "string", "^this "),
		NewDimensionOrderedMetric(LexicographicOrder),
		[]*Aggregator{NewCountAggregator("thing")},
		10,
	)

	res, err := query.ToJSON()
	if err != nil {
		t.Error("got an error", err)
	}
	t.Log(string(res))
}

func TestTimeseries(t *testing.T) {
	origin, _ := time.Parse(time.RFC3339, "2019-02-01T08:00:00.000Z")
	granularity := NewPeriodGranularity(P1M, "America/Los_Angeles", origin)
	query := NewTimeseriesQuery("test", []string{"something"}, []*Aggregator{NewCountAggregator("thing")}, granularity, false)
	res, err := query.ToJSON()
	if err != nil {
		t.Error("got an error", err)
	}
	t.Log(string(res))
}

func TestComplexTimeseries(t *testing.T) {

	// Create granularity
	origin, _ := time.Parse(time.RFC3339, "2020-01-20T00:00:00.00Z")
	granularity := NewPeriodGranularity(P1M, "Europe/Prague", origin)

	// Create context
	qCtx := &QueryContext{QueryID: "activeusers", Priority: 1}

	// Create filters
	baseF := NewBaseFilter(AndExp)
	appIDFilter := NewDimensionFilter("application_id", "123")
	// Events exclusion
	baseEventExclude := NewBaseFilter(NotExp)
	orEventsFilter := NewBaseFilter(OrExp)
	orEventsFilter.AddChildren(
		NewDimensionFilter("event_name", "SessionInitial"),
		NewDimensionFilter("event_name", "PurchaseEvent"),
		NewDimensionFilter("event_name", "Registration"),
		NewDimensionFilter("event_name", "LevelOne"),
		NewDimensionFilter("event_name", "LevelTwo"),
	)
	baseEventExclude.AddChildren(orEventsFilter)
	baseF.AddChildren(baseEventExclude, appIDFilter)

	// Create aggregations
	hyperU := NewHyperUniqueAggregator("unique_devices", "ACTIVE_USERS")
	filt := NewDimensionFilter("event_name", "Purchase")
	agg := NewDoubleSumAggregator("currency_sum", "REV")
	rev := NewFilteredAggregator("REV", filt, agg)

	agg2 := NewDoubleSumAggregator("dimension_count", "Purchases")
	purch := NewFilteredAggregator("Purchases", filt, agg2)

	// Create PostAggregation
	postAgg := NewArithmeticPostAggregation(DivideFunction, "Revenue_Per_User")
	revField := NewFieldAccessPostAggregation("REV")
	activeUsersField := NewHyperUniquePostAggregator("ACTIVE_USERS")
	postAgg.AddChildren(revField, activeUsersField)

	query := NewTimeseriesQuery("events", []string{"2019-02-01T08/2020-01-28T09:14:48"}, []*Aggregator{rev, purch, hyperU}, granularity, false)
	query.AddContext(qCtx)
	query.AddFilter(baseF)
	query.AddPostAggregation(postAgg)
	res, err := query.ToJSON()
	if err != nil {
		t.Error("got an error", err)
	}
	t.Log(string(res))
}

func TestReadMeTimeseries(t *testing.T) {
	// Create granularity
	origin, _ := time.Parse(time.RFC3339, "2020-01-20T00:00:00.00Z")
	granularity := NewPeriodGranularity(P1M, "Europe/Prague", origin)

	// Create context
	qCtx := &QueryContext{
		QueryID:  "activeusers",
		Priority: 1,
		Timeout:  30000,
	}

	// Create filters
	baseFilter := NewBaseFilter(AndExp)
	baseFilter.AddChildren(
		NewDimensionFilter("application_id", "12345"),
		NewDimensionFilter("application_version", "v1.0.1"),
	)

	// Create aggregators
	hyperU := NewHyperUniqueAggregator("unique_devices", "ACTIVE_USERS")

	// Create query
	query := NewTimeseriesQuery(
		"events",
		[]string{"2019-02-01T08/2020-01-28T09:14:48"},
		[]*Aggregator{hyperU},
		granularity,
		false,
	)
	query.AddFilter(baseFilter)
	query.AddContext(qCtx)

	// Create client
	client := NewClient([]string{"http://localhost:8082", "http://localhost:8083"}, "/druid/v2/?pretty")
	res, err := client.Run(query)
	if err != nil {
		fmt.Println(err) // Should give a detailed description of query run, status code and the druid error
	}
	fmt.Println(string(res)) // Return value in bytes to allow caller to manipulate as needed.
}

func TestReadMeTopN(t *testing.T) {
	// Create simple topN query
	query := NewTopNQuery(
		"test",
		[]string{"something"},
		NewPartialExtractionDimension("custom_properties", "this", "string", "^this "),
		NewDimensionOrderedMetric(LexicographicOrder),
		[]*Aggregator{NewCountAggregator("count")},
		10,
	)

	// Create client
	client := NewClient([]string{"http://localhost:8082", "http://localhost:8083"}, "/druid/v2/?pretty")
	res, err := client.Run(query)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(res))
}

func TestReadMeGroupBy(t *testing.T) {
	origin, _ := time.Parse(time.RFC3339, "2020-01-20T00:00:00.00Z")
	granularity := NewPeriodGranularity(P1D, "Europe/Prague", origin)

	// Create query
	query := NewGroupByQuery(
		"some_data",
		granularity,
		[]string{"2020-01-20T00:00:00.000/2020-01-21T00:00:00.000"}, []*Aggregator{NewCountAggregator("count")},
		[]string{"application_id"},
	)

	// Create filters
	filter := NewBaseFilter(OrExp)
	appIDFilter := NewDimensionFilter("application_id", "1347123")
	appIDFilter2 := NewDimensionFilter("application_id", "1010123")
	filter.AddChildren(appIDFilter, appIDFilter2)
	query.AddFilter(filter)

	// Create client
	client := NewClient([]string{"http://localhost:8082", "http://localhost:8083"}, "/druid/v2/?pretty")
	res, err := client.Run(query)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(res))
}
