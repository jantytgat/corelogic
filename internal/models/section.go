package models

import (
	"fmt"
	"strings"
)

type Section struct {
	Name     string    `yaml:"name"`
	Elements []Element `yaml:"elements"`
}

func (s *Section) GetFullName(moduleName string) string {
	return moduleName + "." + s.Name
}

func (s *Section) expandSectionPrefix(expression string) string {
	return strings.ReplaceAll(expression, "prefix", s.Name)
}

func (s *Section) GetFields(moduleName string) (map[string]string, error) {
	output := make(map[string]string)
	var err error

	for _, e := range s.Elements {
		elementOutputName := e.GetFullName(s.GetFullName(moduleName))
		for _, f := range e.Fields {
			outputKey := elementOutputName + "/" + f.Id
			if _, isMapContainsKey := output[outputKey]; isMapContainsKey {
				err = fmt.Errorf("duplicate key in fields: %q", outputKey)
				break
			} else {
				output[outputKey] = s.expandSectionPrefix(f.Data)
			}
		}
	}

	return output, err
}

func (s *Section) GetInstallExpressions(moduleName string, tagFilter []string) (map[string]string, error) {
	output := make(map[string]string)
	var err error

	for _, e := range s.Elements {
		if !e.HasFilteredTag(tagFilter) {
			outputKey := e.GetFullName(s.GetFullName(moduleName))
			var outputValue string
			outputValue, err = e.GetFullyQualifiedExpression(e.Expressions.Install, s.GetFullName(moduleName))

			if err != nil {
				break
			} else {
				if _, isMapContainsKey := output[outputKey]; isMapContainsKey {
					//key exist
					err = fmt.Errorf("duplicate key in section: %q", outputKey)
					break
				} else {
					output[outputKey] = s.expandSectionPrefix(outputValue)
				}
			}
		}
	}

	return output, err
}

func (s *Section) GetUninstallExpressions(moduleName string, tagFilter []string) (map[string]string, error) {
	output := make(map[string]string)
	var err error

	for _, e := range s.Elements {
		if !e.HasFilteredTag(tagFilter) {
			outputKey := e.GetFullName(s.GetFullName(moduleName))
			var outputValue string
			outputValue, err = e.GetFullyQualifiedExpression(e.Expressions.Uninstall, s.GetFullName(moduleName))

			if err != nil {
				break
			} else {
				if _, isMapContainsKey := output[outputKey]; isMapContainsKey {
					//key exist
					err = fmt.Errorf("duplicate key in section: %q", outputKey)
					break
				} else {
					output[outputKey] = s.expandSectionPrefix(outputValue)
				}
			}
		}
	}

	return output, err
}
