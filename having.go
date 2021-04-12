package godruid

type Having struct {
	havingType  string
	filter      *Filter
	aggregation string
	value       string
	havings     []*Having
}

func (h *Having) addHaving(operation OperationType, having *Having) {}
