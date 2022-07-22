package config



// configurationProviderImpl is implementation of IConfigurationProvider.
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

func (c *configurationProviderImpl) GetReloadToken() IChangeToken {
	panic("implement me")
}

func (c *configurationProviderImpl) GetChildKeys(earlierKeys []string, parentPath string) {
	panic("implement me")
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
//	// GetReloadToken returns a `IChangeToken` that can be used to observe when
//	// this configuration is reloaded.
//	GetReloadToken() IChangeToken
//}
