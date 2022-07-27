// Package config is a simplified cut-down version of [ASP.NET] configuration extensions.
//
// Main features:
//
//	- Functionality and behaviour should be identical to [ASP.NET].
//	- Support various configuration sources like Json and Environmental Variables.
//	- Support for hierarchical keys.
//	- Case-insensitive key names.
//	- Hopefully simple and intuitive usage.
//
// Additional features:
//
//	- Ability to inspect configuration and determine the source of values.
//	  This is to support various DevOps and troubleshooting tooling.
//
// Motivation:
//
//	- Mainly for tools which need to inspect and report on existing configuration
//	  of [ASP.NET] applications, e.g. Azure App Services, or .NET apps running in
//	  Kubernetes.
//	- Potentially might be handy for .NET developers using Go to use familiar concepts.
//
// Current configuration providers:
//
//	- Json from user-supplied [[]byte] array. File support coming soon.
//	- Environmental Variables.
//	- Environmental Variables from user-supplied [map[string]string].
//
// Limitations and unimplemented features:
//
//	- Currently read-only versions of everything.
//	- No support for refresh notifications.
//	- No support for many sources like INI files but these may be added later.
//
// See examples for basic and more advanced usage.
//
// [ASP.NET]: https://github.com/dotnet/runtime/tree/release/6.0/src/libraries/Microsoft.Extensions.Configuration/src
package config

import (
	"sort"
	"strings"
)

// keyDelimiter is a hierarchical delimiter for keys.
const keyDelimiter = ":"

// normalizeKey is applied to all keys when adding and querying values.
func normalizeKey(key string) string {
	key = strings.ToLower(key)
	key = strings.Replace(key, "__", keyDelimiter, -1)
	return key
}

// Source is a common interface for types which can provide Config, e.g.
// Json source, or Environmental variables source. This is broadly similar to
// ASP.NET IConfigurationSource and IConfigurationBuilder interfaces.
type Source interface {
	Name() string
	Build() (Config, error)
}

// Config is a simplified cut-down version of ASP.NET IConfiguration interface.
// It provides read-only access to keys and values in the same way ASP.NET does.
// In addition to ASP.NET it gives access to all keys and the source which provided
// this Config object.
//
// The main differences between this type and standard Go maps are:
// 	- Keys are hierarchical, using delimiter [config.keyDelimiter].
// 	- Keys are case-insensitive and are normalised, see below.
//
// Key normalisation happens to provide case-insensitive experience like in ASP.NET:
//   - Converted to lower case.
//   - Double underscore  "__" is converted to [config.keyDelimiter].
type Config interface {
	// Get returns a value for the specified key. The keys are normalised and
	// are not case-sensitive. See [config.normalizeKey] function.
	Get(key string) string
	// TryGet returns a value for the specified key and an indicator whether it exists.
	// The keys are normalised and are not case-sensitive. See [config.normalizeKey] function.
	TryGet(key string, val *string) (found bool)
	// Keys lists all keys in configuration. The list is in alphabetical order.
	// Note: this method is not part of .NET IConfiguration.
	Keys() []string
	// Source returns the [config.Source] which provided this Config.
	Source() Source
}

// newConfigImpl creates an instance of Config interface.
func newConfigImpl(configSource Source, m map[string]string) *configImpl {
	m2 := make(map[string]string, len(m))
	for k, v := range m {
		k = normalizeKey(k)
		m2[k] = v
	}

	return &configImpl{
		configSource: configSource,
		m:            m2,
	}
}

// configImpl implements Config interface.
type configImpl struct {
	m            map[string]string
	configSource Source
}

func (c *configImpl) Source() Source {
	return c.configSource
}

func (c *configImpl) TryGet(key string, val *string) (found bool) {
	key = normalizeKey(key)
	*val, found = c.m[key]
	return found
}

func (c *configImpl) Keys() []string {
	var keys []string
	for k := range c.m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func (c *configImpl) Get(key string) string {
	key = normalizeKey(key)
	return c.m[key]
}

// Builder builds a unified [config.Config] object from multiple Sources.
type Builder interface {
	// AddSource adds a source of configuration. The sources are appended to the
	// end of the existing list. The source further down the list take precedence
	// over sources which are earlier in the list.
	AddSource(source Source)
	// Build builds the [config.RootConfig] object from the current list of sources.
	// Once built, any changes to the sources have no effect. If there are changes
	// in the sources, invoke Build again.
	Build() (RootConfig, error)
}

// NewBuilder create new instance of [config.Builder] implementation.
func NewBuilder() Builder {
	return &builderImpl{}
}

// builderImpl implements [config.Builder] interface.
type builderImpl struct {
	sources []Source
}

func (b *builderImpl) AddSource(source Source) {
	b.sources = append(b.sources, source)
}

func (b *builderImpl) Build() (RootConfig, error) {
	root := newRootConfigImpl()
	for _, source := range b.sources {
		config, err := source.Build()
		if err != nil {
			return nil, err
		}
		root.addConfig(config)
	}

	return root, nil
}

// RootConfig is the main top-level configuration object which aggregates
// configuration from multiple configuration sources and provides and "effective
// configuration".
//
// The sources are maintained in a list, and the sources further down the list take
// precedence over sources which are earlier in the list.
//
// This effectively allows layering and overriding of configuration values e.g.
// default config comes from a JSON file, which can be overridden by Environmental
// variables, which in turn can be overridden by use-supplied values or command line
// arguments.
//
// This type implements [config.Config] interface, and provides additional
// methods to inspect the sources of configuration values.
type RootConfig interface {
	Config
	// GetEntry returns a configuration entry which has information about the
	// source of the value.
	GetEntry(key string) Entry
	// GetEntries returns a list of [config.Entry]. The list is sorted by keys.
	GetEntries() []Entry
}

// newRootConfigImpl creates new instance of [config.RootConfig].
func newRootConfigImpl() *rootConfigImpl {
	return &rootConfigImpl{}
}

// rootConfigImpl implements [config.RootConfig]
type rootConfigImpl struct {
	configs []Config
}

func (c *rootConfigImpl) addConfig(config Config) {
	c.configs = append(c.configs, config)
}

func (c *rootConfigImpl) Get(key string) string {
	return c.GetEntry(key).Value()
}

func (c *rootConfigImpl) TryGet(key string, val *string) (found bool) {
	entry, found := c.tryGetEntry(key)
	*val = entry.Value()
	return found
}

func (c *rootConfigImpl) Keys() []string {
	keysSet := make(map[string]interface{})
	for _, config := range c.configs {
		for _, key := range config.Keys() {
			keysSet[key] = nil
		}
	}

	var keys []string
	for key := range keysSet {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

func (c *rootConfigImpl) Source() Source {
	// TODO(ppanyukov): implement proper source for RootConfig
	return nil
}

func (c *rootConfigImpl) GetEntry(key string) Entry {
	entry, _ := c.tryGetEntry(key)
	return entry
}

func (c *rootConfigImpl) GetEntries() []Entry {
	entrySet := make(map[string]Entry)
	for _, config := range c.configs {
		for _, key := range config.Keys() {
			val := config.Get(key)
			entrySet[key] = newEntryImpl(key, val, config.Source())
		}
	}

	var entries []Entry
	for _, v := range entrySet {
		entries = append(entries, v)
	}

	sort.Slice(entries, func(i int, j int) bool {
		left := entries[i]
		right := entries[j]
		return strings.Compare(left.Key(), right.Key()) == -1
	})

	return entries
}

func (c *rootConfigImpl) tryGetEntry(key string) (result Entry, found bool) {
	for i := len(c.configs) - 1; i >= 0; i-- {
		val := ""
		config := c.configs[i]
		if found := config.TryGet(key, &val); found {
			return newEntryImpl(key, val, config.Source()), true
		}
	}
	return newEntryImpl(key, "", nil), false
}

// Entry is a configuration entry combining key, value and the source of the value.
type Entry interface {
	Key() string
	Value() string
	Source() Source
}

// newEntryImpl creates new instance of [config.Entry]
func newEntryImpl(key string, value string, configSource Source) *configEntryImpl {
	return &configEntryImpl{
		key:          key,
		value:        value,
		configSource: configSource,
	}
}

// configEntryImpl implements [config.Entry] interface.
type configEntryImpl struct {
	key          string
	value        string
	configSource Source
}

func (c *configEntryImpl) Key() string {
	return c.key
}

func (c *configEntryImpl) Value() string {
	return c.value
}

func (c *configEntryImpl) Source() Source {
	return c.configSource
}
