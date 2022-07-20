package interfaces

type SectionReader interface {
	GetFullName(moduleName string) string
	expandSectionPrefix(expression string) string
	GetFields(moduleName string) (map[string]string, error)
	GetInstallExpressions(moduleName string) (map[string]string, error)
	GetUninstallExpressions(moduleName string) (map[string]string, error)
}
