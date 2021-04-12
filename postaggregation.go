package godruid

type PostAggregationType string

type GreatestLeastFunctionName string

type ArithemeticFunction string

type ThetaFunction string

var (
	// PostAggregationType
	ArithmeticPostAggregation            PostAggregationType = "arithmetic"
	FieldAccessPostAggregation           PostAggregationType = "fieldAccess"
	FinalizingFieldAccessPostAggregation PostAggregationType = "finalizingFieldAccess"
	ConstantPostAggregation              PostAggregationType = "constant"
	HyperUniqueAggregation               PostAggregationType = "hyperUniqueCardinality"
	ThetaSketchEstimatePostAggregation   PostAggregationType = "thetaSketchEstimate"
	ThetaSketchSetOpPostAggregation      PostAggregationType = "thetaSketchSetOp"
	// GreatestLeastFunctionName
	DoubleGreatest GreatestLeastFunctionName = "doubleGreatest"
	LongGreatest   GreatestLeastFunctionName = "longGreatest"
	DoubleLeast    GreatestLeastFunctionName = "doubleLeast"
	LongLeast      GreatestLeastFunctionName = "longLeast"
	// ArithemeticFunction
	PlusFunction     ArithemeticFunction = "+"
	MinusFunction    ArithemeticFunction = "-"
	MultiplyFunction ArithemeticFunction = "*"
	DivideFunction   ArithemeticFunction = "/"
	QuotientFunction ArithemeticFunction = "quotient"
	// ThetaFunction
	IntersectTheta ThetaFunction = "INTERSECT"
	UnionTheta     ThetaFunction = "UNION"
	NotTheta       ThetaFunction = "NOT"
)

type PostAggregation struct {
	Type       PostAggregationType `json:"type"`
	Name       string              `json:"name,omitempty"`
	Fields     []*PostAggregation  `json:"fields,omitempty"`
	Field      *PostAggregation    `json:"field,omitempty"`
	FieldNames []string            `json:"fieldNames,omitempty"`
	FieldName  string              `json:"fieldName,omitempty"` // ffs druid
	Fn         ArithemeticFunction `json:"fn,omitempty"`
	Func       ThetaFunction       `json:"func,omitempty"` // ffs again
	Function   string              `json:"function,omitempty"`
}

func (p *PostAggregation) AddChildren(postAggs ...*PostAggregation) {
	p.Fields = append(p.Fields, postAggs...)
}

func NewArithmeticPostAggregation(fn ArithemeticFunction, name string) *PostAggregation {
	return &PostAggregation{
		Type: ArithmeticPostAggregation,
		Fn:   fn,
		Name: name,
	}
}

func NewFieldAccessPostAggregation(fieldName string) *PostAggregation {
	return &PostAggregation{
		Type:      FieldAccessPostAggregation,
		FieldName: fieldName,
	}
}

func NewHyperUniquePostAggregator(fieldName string) *PostAggregation {
	return &PostAggregation{
		Type:      HyperUniqueAggregation,
		FieldName: fieldName,
	}
}

func NewGreatestLeastPostAggregation(functionName GreatestLeastFunctionName, name string, fieldNames []string) *PostAggregation {
	return &PostAggregation{
		Type: PostAggregationType((string)(functionName)),
	}
}

func NewThetaSketchEstimatePostAggregator(metricName, outputName string, fn ThetaFunction, fields ...*PostAggregation) *PostAggregation {
	return &PostAggregation{
		Type: ThetaSketchEstimatePostAggregation,
		Name: outputName,
		Field: &PostAggregation{
			Type:   ThetaSketchSetOpPostAggregation,
			Func:   fn,
			Name:   outputName + ".sketch", // TODO: Not sure how other do this
			Fields: fields,
		},
	}
}
