package models

import (
	"fmt"
	"reflect"
	"testing"
)

func TestElement_GetFullName(t *testing.T) {
	e := Element{
		Name:        "name",
		Tags:        nil,
		Fields:      nil,
		Expressions: Expression{},
	}

	result := e.GetFullName("moduleName")
	want := "moduleName.name"

	if result != want {
		t.Errorf("GetFullName(\"moduleName\") = %s, expected %s", result, want)
	}
}

func TestElement_GetFields(t *testing.T) {
	e := Element{
		Name: "elementName",
		Tags: nil,
		Fields: []Field{
			{
				Id:     "id",
				Data:   "data",
				Prefix: false,
			},
		},
		Expressions: Expression{},
	}

	result, _ := e.GetFields("moduleName")
	want := make(map[string]string)
	want["moduleName.elementName/id"] = "data"

	if reflect.DeepEqual(result, want) != true {
		t.Errorf("GetFields(\"moduleName\") = %s, expected %s", result, want)
	}

}

func TestElement_GetFields2(t *testing.T) {
	e := Element{
		Name: "elementName",
		Tags: nil,
		Fields: []Field{
			{
				Id:     "id",
				Data:   "data",
				Prefix: false,
			},
			{
				Id:     "id",
				Data:   "data",
				Prefix: false,
			},
		},
		Expressions: Expression{},
	}

	_, err := e.GetFields("moduleName")
	if err == nil {
		t.Errorf("Expected error, got %s", err)
	}
}

func TestElement_GetFullyQualifiedExpression(t *testing.T) {
	e := Element{
		Name: "elementName",
		Tags: nil,
		Fields: []Field{
			{
				Id:     "field1",
				Data:   "data1",
				Prefix: false,
			},
			{
				Id:     "field2",
				Data:   "<<otherModule.otherElement/field1>>",
				Prefix: false,
			},
		},
	}
	moduleName := "moduleName"

	var tests = []struct {
		expression string
		moduleName string
		want       string
	}{
		{"install expression", moduleName, "install expression"},
		{"install expression <<field1>>", moduleName, "install expression data1"},
		{"install <<field1>> expression", moduleName, "install data1 expression"},
		{"<<field1>> install expression", moduleName, "data1 install expression"},
		{"install expression <<field1>> <<field2>>", moduleName, "install expression data1 <<otherModule.otherElement/field1>>"},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s, %s", tt.expression, tt.moduleName)
		t.Run(testName, func(t *testing.T) {
			result, _ := e.GetFullyQualifiedExpression(tt.expression, tt.moduleName)
			want := tt.want
			if result != want {
				t.Errorf("GetFullyQualifiedExpression(%s, %s) = %s, expected: %s", tt.expression, tt.moduleName, result, want)
			}
		})
	}

}

func TestElement_ElementHasFilteredTag(t *testing.T) {
	e := Element{
		Name:        "elementName",
		Tags:        nil,
		Fields:      nil,
		Expressions: Expression{},
	}

	var tests = []struct {
		elementTags []string
		filterTags  []string
		want        bool
	}{
		{[]string{"tag1", "tag2"}, []string{""}, false},
		{[]string{"tag1", "tag2"}, []string{"tag1"}, true},
		{[]string{"tag1", "tag2"}, []string{"tag2"}, true},
		{[]string{"tag1", "tag2"}, []string{"tag3"}, false},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s, %s", tt.elementTags, tt.filterTags)
		t.Run(testName, func(t *testing.T) {
			e.Tags = tt.elementTags
			result := e.HasFilteredTag(tt.filterTags)
			want := tt.want

			if result != tt.want {
				t.Errorf("ElementHasFilteredTags(%s) = %v, expected %v", tt.filterTags, result, want)
			}
		})
	}

}
