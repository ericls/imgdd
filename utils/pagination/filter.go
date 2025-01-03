package pagination

type FilterOperator string

const (
	FilterOperatorEq       FilterOperator = "eq"
	FilterOperatorNeq      FilterOperator = "neq"
	FilterOperatorGt       FilterOperator = "gt"
	FilterOperatorGte      FilterOperator = "gte"
	FilterOperatorLt       FilterOperator = "lt"
	FilterOperatorLte      FilterOperator = "lte"
	FilterOperatorContains FilterOperator = "contains"
)

type FilterField struct {
	Name     string
	Operator FilterOperator
	Value    string
}

type Filter struct {
	Fields []FilterField
}

func (f *Filter) AddFilterField(name string, operator FilterOperator, value string) {
	f.Fields = append(f.Fields, FilterField{
		Name:     name,
		Operator: operator,
		Value:    value,
	})
}
