package godruid

import (
	"encoding/json"
	"time"
)

type Period string

var (
	P1D Period = "P1D"
	P1W Period = "P1W"
	P1M Period = "P1M"
)

type Granularity struct {
	GranularityType string    `json:"type"`
	Period          Period    `json:"period"`
	TimeZone        string    `json:"timeZone"`
	Origin          time.Time `json:"origin"`
}

func NewPeriodGranularity(period Period, timeZone string, origin time.Time) Granularity {
	return Granularity{
		GranularityType: "period",
		Period:          period,
		TimeZone:        timeZone,
		Origin:          origin,
	}
}

func NewAllGranularity() Granularity {
	return Granularity{
		GranularityType: "all",
	}
}

func (g Granularity) MarshalJSON() ([]byte, error) {
	type Alias Granularity
	if g.GranularityType == "all" {
		return []byte("\"all\""), nil
	}
	return json.Marshal((Alias)(g))
}
