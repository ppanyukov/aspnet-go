package config

// Equivalent of Microsoft.Extensions.Configuration.Abstractions package.
//
// Defines all sorts of interfaces.

// IChangeToken propagates notifications that a change has occurred.
type IChangeToken interface {
	// HasChanged gets a value that indicates if a change has occurred.
	HasChanged() bool
	// ActiveChangeCallbacks indicates if this token will pro-actively raise callbacks. If `false`,
	// the token consumer must poll `HasChanged` to detect changes.
	ActiveChangeCallbacks() bool
	// RegisterChangeCallback registers for a callback that will be invoked when the entry has changed.
	RegisterChangeCallback(callback func(state interface{}), state interface{})
}

// IConfiguration represents a set of key/value application configuration properties.
type IConfiguration interface {
	Get(key string) (val string)
	Set(key string, val string)
	GetSection(key string) IConfigurationSection
	GetChildren() []IConfigurationSection
	GetReloadToken() IChangeToken
}

// IConfigurationBuilder represents a type used to build application configuration.
type IConfigurationBuilder interface {
	// Properties gets a key/value collection that can be used to share data between the
	// `IConfigurationBuilder` and the registered `IConfigurationSource`s.
	Properties() map[string]interface{}
	// Sources gets the sources used to obtain configuration values.
	Sources() []IConfigurationSource
	// Add adds a new configuration source.
	Add(source IConfigurationSource) IConfigurationBuilder
	// Build builds an `IConfiguration` with keys and values from the set of sources registered in `Sources`.
	Build() IConfigurationRoot
}

// IConfigurationProvider provides configuration key/values for an application.
//
// This interface is equivalent to IConfigurationProvider interface in .NET.
// See: https://github.com/dotnet/runtime/blob/release/6.0/src/libraries/Microsoft.Extensions.Configuration.Abstractions/src/IConfigurationProvider.cs
type IConfigurationProvider interface {
	// Get gets a configuration value for the specified key. Returns empty string if not found.
	Get(key string) (value string)
	// TryGet tries to get a configuration value for the specified key.
	// TODO: convert TryGet to Go style?
	TryGet(key string, val *string) (found bool)
	// Set sets a configuration value for the specified key.
	Set(key string, val string)
	// GetReloadToken returns a change token if this provider supports change tracking, null otherwise.
	// TODO: do we need GetReloadToken?
	GetReloadToken() IChangeToken
	// Load loads configuration values from the source represented by this `IConfigurationProvider`.
	Load()
	// GetChildKeys returns the immediate descendant configuration keys for a given parent path based on this
	// `IConfigurationProvider` data and the set of keys returned by all the preceding `IConfigurationProvider`.
	// TODO: do we need GetChildKeys? What does it even do??
	GetChildKeys(earlierKeys []string, parentPath string)
}

// IConfigurationRoot represents the root of an `IConfiguration` hierarchy.
type IConfigurationRoot interface {
	IConfiguration
	// Reload Forces the configuration values to be reloaded from the underlying `IConfigurationProvider`s.
	Reload()
	// Providers gets the list of `IConfigurationProvider`s in the configuration.
	Providers() []IConfigurationProvider
}

// IConfigurationSection represents a section of application configuration values.
type IConfigurationSection interface {
	IConfiguration
	Key() string
	Path() string
	Value() string
}

// IConfigurationSource represents a source of configuration key/values for an application.
type IConfigurationSource interface {
	// Build builds the `IConfigurationProvider` for this source.
	Build(builder IConfigurationBuilder) IConfigurationProvider
}
