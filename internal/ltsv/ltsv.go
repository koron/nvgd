// Package ltsv provides LTSV utility functioons.
package ltsv

// Set is a set of LTSV values in a line.
type Set struct {
	Properties []Property
	Index      map[string][]int
}

// Property is a pair of label and value.
type Property struct {
	Label string
	Value string
}

func NewSet() *Set {
	return &Set{
		Index: make(map[string][]int),
	}
}

// Put puts a property to the set.
func (s *Set) Put(label, value string) {
	n := len(s.Properties)
	s.Properties = append(s.Properties, Property{Label: label, Value: value})
	s.Index[label] = append(s.Index[label], n)
}

func (s *Set) PutProperties(props []Property) {
	for _, p := range props {
		s.Put(p.Label, p.Value)
	}
}

// Get gets values for the label.
func (s *Set) Get(label string) []string {
	indexes, ok := s.Index[label]
	if !ok {
		return nil
	}
	list := make([]string, len(indexes))
	for i, n := range indexes {
		list[i] = s.Properties[n].Value
	}
	return list
}

// GetFirst gets a first value for the label.
func (s *Set) GetFirst(label string) string {
	indexes, ok := s.Index[label]
	if !ok || len(indexes) == 0 {
		return ""
	}
	return s.Properties[indexes[0]].Value
}

// Empty checks set is empty or not.
func (s *Set) Empty() bool {
	return len(s.Properties) == 0
}
