package interfaces

type ModuleReader interface {
	GetFullModuleName(packageName string) string
	GetFields(packageName string) (map[string]string, error)
	GetInstallExpressions(packageName string, tagFilter []string) (map[string]string, error)
	GetUninstallExpressions(packageName string, tagFilter []string) (map[string]string, error)
	AppendData(source map[string]string, destination map[string]string) (map[string]string, error)
	HasFilteredTag(tagFilter []string) bool
}
