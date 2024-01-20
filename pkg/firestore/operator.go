package firestore

type operator struct {
	Equal         string
	NotEqual      string
	LowerThan     string
	LowerOrEqual  string
	HigherThan    string
	HigherOrEqual string
	In            string
	NotIn         string
}

var Operator = operator{
	Equal:         "==",
	NotEqual:      "!=",
	LowerThan:     "<",
	LowerOrEqual:  "<=",
	HigherThan:    ">",
	HigherOrEqual: ">=",
	In:            "in",
	NotIn:         "not-in",
}
