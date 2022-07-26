package config

import (
	"fmt"
	"strings"
)

// envVarsLoader is an implementation of .NET EnvironmentVariablesConfigurationProvider.
//
// It does it in the same way as ASP.NET EnvironmentVariablesConfigurationProvider.
// The actual code was broadly taken as-is from ASP.NET.
//
// The tests were ported from ASP.NET too, so the behaviour is identical.
//
// See: https://github.com/dotnet/runtime/blob/release/6.0/src/libraries/Microsoft.Extensions.Configuration.EnvironmentVariables/src/EnvironmentVariablesConfigurationProvider.cs
type envVarsLoader struct {
	prefix string
}

// newEnvVarsLoader creates a loader of config from env variables.
func newEnvVarsLoader(prefix string) *envVarsLoader {
	e := &envVarsLoader{
		// Don't use normalize, the "__" should not be replaced in prefix, just lower case.
		// This is how ASP.NET behaves, there is even a test for this.
		prefix: strings.ToLower(prefix),
	}

	return e
}

// Env vars with these prefixes have special meaning in .NET.
// These are used
var (
	envMySqlServerPrefix    = normalizeKey("MYSQLCONNSTR_")
	envSqlAzureServerPrefix = normalizeKey("SQLAZURECONNSTR_")
	envSqlServerPrefix      = normalizeKey("SQLCONNSTR_")
	envCustomPrefix         = normalizeKey("CUSTOMCONNSTR_")
)

func (e *envVarsLoader) Load(envVars map[string]string) map[string]string {
	m := make(map[string]string, len(envVars))

	// NOTE: the weird logic around SQL is taken as-is from ASP.NET codebase.
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

	return m
}

func (e *envVarsLoader) addIfPrefixed(m map[string]string, key string, val string) {
	if strings.HasPrefix(key, e.prefix) {
		key = strings.TrimPrefix(key, e.prefix)
		key = strings.TrimPrefix(key, keyDelimiter)
		m[key] = val
	}
}

func (e *envVarsLoader) trimPrefix(key string, prefix string) string {
	key = strings.TrimPrefix(key, prefix)
	key = strings.TrimPrefix(key, keyDelimiter)
	return key
}
