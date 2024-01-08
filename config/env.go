package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

const (
	baseUrlEnvName = "BASE_URL"
	baseUrlDefault = "http://upstream:80/_/api"
)

var (
	exampleSvcSuffix = "integrations/example-svc"

	BaseServiceNotSpecifiedErr = errors.New("base service was not specified")
)

type testsEnvConfig struct {
	baseUrl     string
	baseService *string
	suffix      string
}

func NewTestsEnvConfig() TestsConfig {
	baseUrl, present := os.LookupEnv(baseUrlEnvName)
	if !present {
		baseUrl = baseUrlDefault
	}
	return testsEnvConfig{
		baseUrl: baseUrl,
	}
}

func (cfg testsEnvConfig) FromExampleSvc() TestsConfig {
	cfg.baseService = &exampleSvcSuffix
	return cfg
}

func (cfg testsEnvConfig) WithSuffix(suffix string) TestsConfig {
	cfg.suffix = suffix
	return cfg
}

func (cfg testsEnvConfig) Build() (*string, error) {
	if cfg.baseService == nil {
		return nil, BaseServiceNotSpecifiedErr
	}

	url := fmt.Sprintf("%s/%s", cfg.baseUrl, *cfg.baseService)
	if len(cfg.suffix) > 0 {
		url = fmt.Sprintf("%s/%s", url, cfg.suffix)
	}

	return &url, nil
}

func (cfg testsEnvConfig) MustBuild() string {
	url, err := cfg.Build()
	if err != nil {
		panic(errors.Wrap(err, "failed to build an url"))
	}

	return *url
}
