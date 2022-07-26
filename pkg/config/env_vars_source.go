package config

import (
	"fmt"
	"os"
	"strings"
)

// NewEnvVarsSource creates configuration source from environmental variables
// in the same way ASP.NET does.
func NewEnvVarsSource(prefix string) *envVarsSource {
	envVars := make(map[string]string)
	for _, val := range os.Environ() {
		fields := strings.SplitN(val, "=", 2)
		envVars[fields[0]] = fields[1]
	}

	return &envVarsSource{
		name:   fmt.Sprintf("envVarsSource Prefix: '%s'", prefix),
		prefix: prefix,
		m:      envVars,
	}
}

// NewEnvVarsMapSource is same as NewEnvVarsSource except the
// values come from the user-supplied map.
func NewEnvVarsMapSource(prefix string, m map[string]string) *envVarsSource {
	return &envVarsSource{
		name:   fmt.Sprintf("envVarsMapSource Prefix: '%s'", prefix),
		prefix: prefix,
		m:      m,
	}
}

type envVarsSource struct {
	name   string
	prefix string
	m      map[string]string
}

// WithName sets the name of this source and returns itself.
func (s *envVarsSource) WithName(name string) *envVarsSource {
	s.name = name
	return s
}

// Name is the name of this source. Part of [config.Source] interface.
func (s *envVarsSource) Name() string {
	return s.name
}

func (s *envVarsSource) Build() (Config, error) {
	m := newEnvVarsLoader(s.prefix).Load(s.m)
	return newConfigImpl(s, m), nil
}
