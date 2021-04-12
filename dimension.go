package godruid

type DimensionType string
type ExtractionFnType string

var (
	// Dimension
	DefaultDimensionType    DimensionType = "default"
	ExtractionDimensionType DimensionType = "extraction"
	// Extraction
	PartialExtractionType    ExtractionFnType = "partial"
	TimeFormatExtractionType ExtractionFnType = "timeFormat"
)

type Dimension struct {
	Type         DimensionType `json:"type"`
	Dimension    string        `json:"dimension"`
	OutputName   string        `json:"outputName"`
	OutputType   string        `json:"outputType,omitempty"` // TODO: Enum
	ExtractionFn *ExtractionFn `json:"extractionFn,omitempty"`
}

type ExtractionFn struct {
	Type       ExtractionFnType `json:"type"`
	Expression string           `json:"expr,omitempty"`
	Format     string           `json:"format,omitempty"`
	TimeZone   string           `json:"timeZone,omitempty"`
	Locale     string           `json:"locale,omitempty"`
}

func NewDefaultDimension(dimensionName, outputName, outputType string) *Dimension {
	return &Dimension{
		Type:       DefaultDimensionType,
		Dimension:  dimensionName,
		OutputName: outputName,
		OutputType: outputType,
	}
}

func NewPartialExtractionDimension(dimensionName, outputName, outputType, expression string) *Dimension {
	extraction := &ExtractionFn{
		Type:       PartialExtractionType,
		Expression: expression,
	}
	return &Dimension{
		Type:         ExtractionDimensionType,
		Dimension:    dimensionName,
		OutputName:   outputName,
		OutputType:   outputType,
		ExtractionFn: extraction,
	}
}

// TODO: Add the rest of the extraction functions
// https://druid.apache.org/docs/latest/querying/dimensionspecs.html
