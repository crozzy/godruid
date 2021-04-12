package godruid

import (
	"encoding/json"
	"errors"
	"time"
)

type OperationType string

var (
	Selector    OperationType = "selector"
	AndExp      OperationType = "and"
	OrExp       OperationType = "or"
	NotExp      OperationType = "not"
	InExp       OperationType = "in"
	BoundFilter OperationType = "bound"
	Regex       OperationType = "regex"

	ErrFilterCannotHaveChildren = errors.New("cannot append filter to selector type")
	ErrFilterCannotHaveSibling  = errors.New("non-selector filter type cannot have sibling")
)

type Filter struct {
	FilterType   OperationType `json:"type,omitempty"`
	Dimension    string        `json:"dimension,omitempty"`
	Value        string        `json:"value,omitempty"`
	Values       []string      `json:"values,omitempty"` // Thanks Druid
	Pattern      string        `json:"pattern,omitempty"`
	Filters      []*Filter     `json:"fields,omitempty"`
	Filter       *Filter       `json:"field,omitempty"`
	Upper        string        `json:"upper,omitempty"`
	Lower        string        `json:"lower,omitempty"`
	ExtractionFn *ExtractionFn `json:"extractionFn,omitmpty"`
}

func NewBaseFilter(filterType OperationType) *Filter {
	return &Filter{
		FilterType: filterType,
		Filters:    []*Filter{},
	}
}

func NewDimensionFilter(dimension, value string) *Filter {
	return &Filter{
		FilterType: Selector,
		Dimension:  dimension,
		Value:      value,
	}
}

func NewRegexFilter(dimension, pattern string) *Filter {
	return &Filter{
		FilterType: Regex,
		Pattern:    pattern,
	}
}

func NewInFilter(dimension string, values []string) *Filter {
	return &Filter{
		FilterType: InExp,
		Dimension:  dimension,
		Values:     values,
	}
}

func NewTimeBoundFilter(dimension string, start, end time.Time, timezone, timeFormat string) *Filter {
	return &Filter{
		FilterType: BoundFilter,
		Upper:      end.Format(timeFormat),
		Lower:      start.Format(timeFormat),
		Dimension:  dimension,
		ExtractionFn: &ExtractionFn{
			Type:     TimeFormatExtractionType,
			Format:   timeFormat,
			TimeZone: timezone,
		},
	}
}

func (f *Filter) AddChildren(newFs ...*Filter) error {
	if f.FilterType == Selector {
		return ErrFilterCannotHaveChildren
	}
	f.Filters = append(f.Filters, newFs...)
	return nil
}

// MarshalJSON will check the length of the child filters as Druid hates 1 item
// arrays and for Not filters they are illegal it needs to be an object.
func (f Filter) MarshalJSON() ([]byte, error) {
	type Alias Filter
	if len(f.Filters) == 1 {
		f.Filter = f.Filters[0]
		f.Filters = nil
	}
	return json.Marshal((Alias)(f))
}
