package models

import (
	"testing"
)

func TestFieldGetModuleName(t *testing.T) {
	f := Field{
		Id:     "id",
		Data:   "data",
		Prefix: false,
	}
	result := f.GetFullName("moduleName")
	want := "moduleName/id"
	if result != want {
		t.Errorf("GetModuleName(\"modulename\") = %s; want %s", result, want)
	}
}

//TODO: add fuzzing?
