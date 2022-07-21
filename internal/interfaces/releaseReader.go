package interfaces

type ReleaseReader interface {
	GetVersionAsString() string
	GetSemanticVersion() string
}
