package config

import (
	"fmt"
	"os"
	"strings"
)

// Env vars with thesse prefixes have special meaning in .NET.
var (
	envMySqlServerPrefix    = normalizeKey("MYSQLCONNSTR_")
	envSqlAzureServerPrefix = normalizeKey("SQLAZURECONNSTR_")
	envSqlServerPrefix      = normalizeKey("SQLCONNSTR_")
	envCustomPrefix         = normalizeKey("CUSTOMCONNSTR_")
)

// envVarsConfigProvider is an implementation of .NET EnvironmentVariablesConfigurationProvider.
// See: https://github.com/dotnet/runtime/blob/release/6.0/src/libraries/Microsoft.Extensions.Configuration.EnvironmentVariables/src/EnvironmentVariablesConfigurationProvider.cs
type envVarsConfigProvider struct {
	configurationProviderImpl
	prefix     string
	prefixOrig string
}

func newEnvVars(prefix string) *envVarsConfigProvider {
	e := &envVarsConfigProvider{
		configurationProviderImpl: *newConfigurationProviderImpl(),
		// Don't use normalize, the "__" should not be replaced in prefix, just lower case.
		// This is how ASP.NET behaves, there is even a test for this.
		prefix:     strings.ToLower(prefix),
		prefixOrig: prefix,
	}

	return e
}

func (e *envVarsConfigProvider) Load() {
	varMap := make(map[string]string)
	for _, val := range os.Environ() {
		fields := strings.SplitN(val, "=", 2)
		varMap[fields[0]] = fields[1]
	}
	e.load(varMap)
}

func (e *envVarsConfigProvider) load(envVars map[string]string) {
	m := make(map[string]string, len(envVars))

	for k, v := range envVars {
		prefix := ""
		provider := ""

		k = normalizeKey(k)

		if strings.HasPrefix(k, envMySqlServerPrefix) {
			prefix = envMySqlServerPrefix
			provider = "MySql.Data.MySqlClient"
		} else if strings.HasPrefix(k, envSqlAzureServerPrefix) {
			prefix = envSqlAzureServerPrefix
			provider = "System.Data.SqlClient"
		} else if strings.HasPrefix(k, envSqlServerPrefix) {
			prefix = envSqlServerPrefix
			provider = "System.Data.SqlClient"
		} else if strings.HasPrefix(k, envCustomPrefix) {
			prefix = envCustomPrefix
		} else {
			e.addIfPrefixed(m, k, v)
			continue
		}

		k = e.trimPrefix(k, prefix)
		k2 := normalizeKey(fmt.Sprintf("ConnectionStrings:%s", k))
		e.addIfPrefixed(m, k2, v)

		if provider != "" {
			k3 := normalizeKey(fmt.Sprintf("ConnectionStrings:%s_ProviderName", k))
			e.addIfPrefixed(m, k3, provider)
		}
	}

	e.m = m
}

func (e *envVarsConfigProvider) addIfPrefixed(m map[string]string, key string, val string) {
	if strings.HasPrefix(key, e.prefix) {
		key = strings.TrimPrefix(key, e.prefix)
		key = strings.TrimPrefix(key, keyDelimiter)
		m[key] = val
	}
}

func (e *envVarsConfigProvider) trimPrefix(key string, prefix string) string {
	key = strings.TrimPrefix(key, prefix)
	key = strings.TrimPrefix(key, keyDelimiter)
	return key
}

func (e *envVarsConfigProvider) String() string {
	return fmt.Sprintf("EnvironmentVariablesConfigurationProvider Prefix: '%s'", e.prefixOrig)
}
