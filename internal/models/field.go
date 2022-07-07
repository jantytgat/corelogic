package models

type Field struct {
	Id     string `yaml:"id"`
	Data   string `yaml:"data"`
	Prefix bool   `yaml:"prefix"`
}

func (f *Field) GetFullName(moduleName string) string {
	return moduleName + "/" + f.Id
}
