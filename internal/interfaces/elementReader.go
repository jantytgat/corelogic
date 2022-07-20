package interfaces

type ElementReader interface {
	GetFullName(moduleName string) string
	GetFields(moduleName string) (map[string]string, error)
	GetFullyQualifiedExpression(expression string, moduleName string) (string, error)
	HasFilteredTag(tagFilter []string) bool
}
