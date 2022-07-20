package interfaces

import "github.com/jantytgat/corelogic/internal/models"

type FrameworkReader interface {
	GetPrefixMap() map[string]string
	GetPrefixWithVersion(sectionName string) string
	GetFields() (map[string]string, error)
	GetExpressions(kind string, tagFilter []string) (map[string]string, error)
	SortPrefixes(prefixes []models.Prefix)
}
