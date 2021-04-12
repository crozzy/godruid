package godruid

type QueryContext struct {
	Timeout                      int64  `json:"timeout,omitempty"`
	MaxScatterGatherBytes        int64  `json:"maxScatterGatherBytes,omitempty"`
	Priority                     int    `json:"priority,omitempty"`
	QueryID                      string `json:"queryId,omitempty"`
	UseCache                     bool   `json:"useCache,omitempty"`
	PopulateCache                bool   `json:"populateCache,omitempty"`
	BySegement                   bool   `json:"bySegment,omitempty"`
	Finalize                     bool   `json:"finalize,omitempty"`
	ChunkPeriod                  string `json:"chunkPeriod,omitempty"`
	SerializeDateTimeAsLong      bool   `json:"serializeDateTimeAsLong,omitempty"`
	SerializeDateTimeAsLonginner bool   `json:"serializeDateTimeAsLonginner,omitempty"`
	// TopN specific
	MinTopNThreshold int64 `json:"minTopNThreshold,omitempty"`
	// Timeseries specific
	SkipEmptyBuckets bool `json:"skipEmptyBuckets"`
	// GroupBy specific
	// TODO
}
