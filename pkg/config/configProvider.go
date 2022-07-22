package config

// ConfigurationProvider provides configuration key/values for an application.
//
// This interface is equivalent to IConfigurationProvider interface in .NET.
// See: https://github.com/dotnet/runtime/blob/release/6.0/src/libraries/Microsoft.Extensions.Configuration.Abstractions/src/IConfigurationProvider.cs
type ConfigurationProvider interface {
	// Get gets a configuration value for the specified key. Returns empty string if not found.
	Get(key string) (value string)
	// TryGet tries to get a configuration value for the specified key.
	TryGet(key string, val *string) (found bool)
	// Set sets a configuration value for the specified key.
	Set(key string, val string)
	// GetReloadToken returns a change token if this provider supports change tracking, null otherwise.
	// TODO: do we need GetReloadToken?
	GetReloadToken() ChangeToken
	// Load loads configuration values from the source represented by this `ConfigurationProvider`.
	Load()
	// GetChildKeys returns the immediate descendant configuration keys for a given parent path based on this
	// `ConfigurationProvider` data and the set of keys returned by all the preceding `ConfigurationProvider`.
	// TODO: do we need GetChildKeys? What does it even do??
	GetChildKeys(earlierKeys []string, parentPath string)
}

// configurationProviderImpl is implementation of ConfigurationProvider.
// Include it as part of
type configurationProviderImpl struct {
	m map[string]string
}

// newConfigurationProviderImpl creates a new instance of configurationProviderImpl
func newConfigurationProviderImpl() *configurationProviderImpl {
	res := &configurationProviderImpl{
		m: make(map[string]string),
	}

	return res
}

func (c *configurationProviderImpl) Get(key string) (value string) {
	key = normalizeKey(key)
	return c.m[key]
}

func (c *configurationProviderImpl) TryGet(key string, val *string) (found bool) {
	key = normalizeKey(key)
	*val, found = c.m[key]
	return found
}

func (c *configurationProviderImpl) Set(key string, val string) {
	key = normalizeKey(key)
	c.m[key] = val
}

func (c *configurationProviderImpl) Load() {
	// by default this does nothing, override in inheritors.
}

func (c *configurationProviderImpl) GetReloadToken() ChangeToken {
	panic("implement me")
}

func (c *configurationProviderImpl) GetChildKeys(earlierKeys []string, parentPath string) {
	panic("implement me")
}

type ChangeToken interface {
	// HasChanged gets a value that indicates if a change has occurred.
	HasChanged() bool
	// ActiveChangeCallbacks indicates if this token will pro-actively raise callbacks. If `false`,
	// the token consumer must poll `HasChanged` to detect changes.
	ActiveChangeCallbacks() bool
	// RegisterChangeCallback registers for a callback that will be invoked when the entry has changed.
	RegisterChangeCallback(callback func(state interface{}), state interface{})
}

//// ConfigurationSection represents a section of application configuration values.
//type ConfigurationSection interface {
//	Configuration
//	// Key gets the key this section occupies in its parent.
//	Key() string
//	// Path gets the full path to this section within the `Configuration`.
//	Path() string
//	// Value gets the section value.
//	Value() string
//}
//
//// Configuration represents a set of key/value application configuration properties.
//type Configuration interface {
//	// Get gets the configuration value.
//	Get(key string) (value string, found bool)
//	// GetSection gets a configuration sub-section with the specified key.
//	GetSection(key string) ConfigurationSection
//	// GetChildren gets the immediate descendant configuration sub-sections.
//	GetChildren() []ConfigurationSection
//	// GetReloadToken returns a `ChangeToken` that can be used to observe when
//	// this configuration is reloaded.
//	GetReloadToken() ChangeToken
//}
