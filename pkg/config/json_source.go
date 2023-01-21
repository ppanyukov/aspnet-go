package config

import (
	"bytes"

	"github.com/pkg/errors"
)

// NewJsonSource creates configuration source for Json which implements [config.Source].
func NewJsonSource(json []byte) *JsonSource {
	return &JsonSource{
		json: json,
		name: "JsonSource",
	}
}

// JsonSource implements [config.Source] interface.
type JsonSource struct {
	json []byte
	name string
}

// WithName sets the name of this source and returns itself.
func (s *JsonSource) WithName(name string) *JsonSource {
	s.name = name
	return s
}

// Name is the name of this source. Part of [config.Source] interface.
func (s *JsonSource) Name() string {
	return s.name
}

// Build builds Config. Part of [config.Source] interface.
func (s *JsonSource) Build() (Config, error) {
	parser := newJsonLoader()
	r := bytes.NewBuffer(s.json)
	m, err := parser.Load(r)
	if err != nil {
		return nil, errors.Errorf("JsonSource: %s: %v", s.name, err)
	}

	return newConfigImpl(s, m), nil
}
