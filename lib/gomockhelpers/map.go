package gomockhelpers

import (
	"fmt"

	"go.uber.org/mock/gomock"
)

func Map(m map[string]interface{}) *mapMatcher {
	return &mapMatcher{m}
}

type mapMatcher struct {
	m map[string]interface{}
}

func (matcher *mapMatcher) Matches(x any) bool {
	m, ok := x.(map[string]interface{})
	if !ok {
		return false
	}

	for k, v := range matcher.m {
		if v == gomock.Any() {
			continue
		}

		if m[k] != v {
			return false
		}
	}

	return true
}

func (matcher *mapMatcher) String() string {
	return fmt.Sprintf("%v", matcher.m)
}
