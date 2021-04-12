# godruid

godruid exposes a simple API to create, execute, and analyze [Druid](http://druid.io/) queries. Build from the standpoint of "give the caller as much control as possible", the client and the query creation are seperate so the caller can replace either should they wish. There is nothing to stop users using the same client in different goroutines should they wish, the client will not open any connection until `Run()` is called on a query. All queries implement the `ToJSON()` method to allow the caller access to the query JSON should they need it.

The instanciating functions for creating queries take the required fields and helper methods for adding the optional fields, hence filters are added after query creation to make the API as clean as possible.

There is minimal hand-holding when it comes to filter, aggregation or postaggregation object, you can easily create a monster that Druid will reject at the moment, it's still to be decided if more guardrails will be added.


# examples

The following exampes show how to execute and analyze the results of three types of queries: timeseries, topN, and groupby. We will use these queries to ask simple questions about twitter's public data set.

## timeseries

```go
    // Create granularity
	origin, _ := time.Parse(time.RFC3339, "2020-01-20T00:00:00.00Z")
	granularity := NewPeriodGranularity(P1M, "Europe/Prague", origin)

    // Create context
    qCtx := &QueryContext{
        QueryID: "activeusers",
        Priority: 1,
        Timeout: 30000,
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
```

## topN

```go
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
```

## groupBy
```go
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
    appIDFilter := NewDimensionFilter("application_id", "134700")
    appIDFilter2 := NewDimensionFilter("application_id", "1012320")
    filter.AddChildren(appIDFilter, appIDFilter2)
    query.AddFilter(filter)

    // Create client
    client := NewClient([]string{"http://localhost:8082", "http://localhost:8083"}, "/druid/v2/?pretty")
    res, err := client.Run(query)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(string(res))

```
