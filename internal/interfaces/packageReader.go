package interfaces

type PackageReader interface {
	GetFields() (map[string]string, error)
	GetInstallExpressions(tagFilter []string) (map[string]string, error)
	GetUninstallExpressions(tagFilter []string) (map[string]string, error)
	AppendData(source map[string]string, destination map[string]string) (map[string]string, error)
}
