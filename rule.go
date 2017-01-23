package typeconv

// Rule represents type conversion rule.
//
// from -> to -> priority
type Rule struct {
	next  int
	rules map[string]map[string]int
}

// Add adds type conversion rule.
func (r *Rule) Add(from, to string) {
	if len(r.rules) == 0 {
		r.rules = make(map[string]map[string]int)
	}
	if _, ok := r.rules[from]; !ok {
		r.rules[from] = make(map[string]int)
	}
	r.rules[from][to] = r.next
	r.next--
}

// ConvertibleTo reports whether a "from" type is convertible to a "to" type.
func (r *Rule) ConvertibleTo(from, to string) (priority int, ok bool) {
	if r.rules == nil {
		return 0, false
	}
	if _, ok := r.rules[from]; !ok {
		return 0, false
	}
	priority, ok = r.rules[from][to]
	return priority, ok
}

// DefaultRule holds default type conversion rules whose conversion are safe.
var DefaultRule = &Rule{}

func init() {
	rules := []struct {
		from string
		to   string
	}{
		{"int8", "int16"},
		{"int8", "int32"},
		{"int8", "int"},
		{"int8", "int64"},
		{"int8", "float32"},
		{"int8", "float64"},

		{"int16", "int32"},
		{"int16", "int"},
		{"int16", "int64"},
		{"int16", "float32"},
		{"int16", "float64"},

		{"int32", "int"},
		{"int32", "int64"},
		{"int32", "float32"},
		{"int32", "float64"},

		{"int", "int64"},
		{"int", "float32"},
		{"int", "float64"},

		{"int64", "float32"},
		{"int64", "float64"},

		{"float32", "float64"},
	}
	for _, r := range rules {
		DefaultRule.Add(r.from, r.to)
	}
}
