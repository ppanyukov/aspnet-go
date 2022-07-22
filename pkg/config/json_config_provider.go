package config

import (
	"io"
)

// jsonConfigProvider implements .NET equivalent of JsonConfigurationProvider.
//
// See: https://github.com/dotnet/runtime/blob/release/6.0/src/libraries/Microsoft.Extensions.Configuration.Json/src/JsonConfigurationProvider.cs
//
// TODO: flesh out jsonConfigProvider, it just implements parser at the moment
type jsonConfigProvider struct {
	configurationProviderImpl
}

func newJsonConfigProvider() *jsonConfigProvider {
	return &jsonConfigProvider{
		configurationProviderImpl: *newConfigurationProviderImpl(),
	}
}

func (j *jsonConfigProvider) Load(r io.Reader) error {
	parser := newJsonConfigFileParser()
	result, err := parser.Parse(r)
	if err != nil {
		return err
	}
	j.m = result
	return nil
}
