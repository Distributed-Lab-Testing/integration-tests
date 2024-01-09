package config

type TestsConfig interface {
	FromExampleSvc() TestsConfig
	WithSuffix(string string) TestsConfig

	Build() (*string, error)
	MustBuild() string
}
