package config

import (
	"bytes"

	"github.com/pkg/errors"
)

// NewJsonSource creates configuration source for Json which implements [config.Source].
func NewJsonSource(json []byte) *jsonSource {
	return &jsonSource{
		json: json,
		name: "jsonSource",
	}
}

// jsonSource implements [config.Source] interface.
type jsonSource struct {
	json []byte
	name string
}

// WithName sets the name of this source and returns itself.
func (s *jsonSource) WithName(name string) *jsonSource {
	s.name = name
	return s
}

// Name is the name of this source. Part of [config.Source] interface.
func (s *jsonSource) Name() string {
	return s.name
}

// Build builds Config. Part of [config.Source] interface.
func (s *jsonSource) Build() (Config, error) {
	parser := newJsonLoader()
	r := bytes.NewBuffer(s.json)
	m, err := parser.Load(r)
	if err != nil {
		return nil, errors.Errorf("%s: %v", s.name, err)
	}

	return newConfigImpl(s, m), nil
}
