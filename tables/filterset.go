package tables

import "github.com/zgljl2012/go-orm"

type filterSet struct {
}

// Filter with paramters
func (f *filterSet) Filter(parameters ...*orm.QueryParameter) orm.FilterSet {
	return nil
}

// OrderBy specify ordering fields, plus means ASC, minus(-) means DESC
func (f *filterSet) OrderBy([]string) orm.FilterSet {
	return nil
}

// Limit rows
func (f *filterSet) Limit(int) orm.FilterSet {
	return nil
}

// All return all rows
func (f *filterSet) All() []interface{} {
	return nil
}
