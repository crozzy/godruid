package godruid

type LimitSpec struct {
	LimitType string   `json:"type"`
	Limit     int64    `json:"limit"`
	Columns   []string `json:"columns"`
}

func NewLimitSpec(limit int64, columns []string) *LimitSpec {
	return &LimitSpec{
		LimitType: "default",
		Limit:     limit,
		Columns:   columns,
	}
}
