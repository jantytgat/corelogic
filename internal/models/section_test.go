package models

import "testing"

func TestSection_GetFullName(t *testing.T) {
	s := Section{
		Name:     "sectionName",
		Elements: nil,
	}
	moduleName := "moduleName"
	result := s.GetFullName(moduleName)
	want := "moduleName.sectionName"

	if result != want {
		t.Errorf("GetFullName(%s) = %s, expected: %s", moduleName, result, want)
	}
}
