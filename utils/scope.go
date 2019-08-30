package utils

import (
	mapset "github.com/deckarep/golang-set"
)

// ScopeIsSubset 是包含
func ScopeIsSubset(input, source []string) bool {
	inputSet := mapset.NewSet()
	for _, i := range input {
		inputSet.Add(i)
	}
	sourceSet := mapset.NewSet()
	for _, s := range source {
		sourceSet.Add(s)
	}
	return inputSet.IsSubset(sourceSet)
}
