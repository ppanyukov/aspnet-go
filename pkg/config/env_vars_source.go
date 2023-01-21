package config

import (
	"fmt"
	"os"
	"strings"
)

// NewEnvVarsSource creates configuration source from environmental variables
// in the same way ASP.NET does.
func NewEnvVarsSource(prefix string) *EnvVarsSource {
	envVars := make(map[string]string)
	for _, val := range os.Environ() {
		fields := strings.SplitN(val, "=", 2)
		envVars[fields[0]] = fields[1]
	}

	return &EnvVarsSource{
		name:   fmt.Sprintf("EnvVarsSource Prefix: '%s'", prefix),
		prefix: prefix,
		m:      envVars,
	}
}

// NewEnvVarsMapSource is same as NewEnvVarsSource except the
// values come from the user-supplied map.
func NewEnvVarsMapSource(prefix string, m map[string]string) *EnvVarsSource {
	return &EnvVarsSource{
		name:   fmt.Sprintf("envVarsMapSource Prefix: '%s'", prefix),
		prefix: prefix,
		m:      m,
	}
}

type EnvVarsSource struct {
	name   string
	prefix string
	m      map[string]string
}

// WithName sets the name of this source and returns itself.
func (s *EnvVarsSource) WithName(name string) *EnvVarsSource {
	s.name = name
	return s
}

// Name is the name of this source. Part of [config.Source] interface.
func (s *EnvVarsSource) Name() string {
	return s.name
}

func (s *EnvVarsSource) Build() (Config, error) {
	m := newEnvVarsLoader(s.prefix).Load(s.m)
	return newConfigImpl(s, m), nil
}
