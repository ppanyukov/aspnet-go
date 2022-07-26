package config

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// These tests are a port of EnvironmentVariablesTest.cs
// https://github.com/dotnet/runtime/blob/release/6.0/src/libraries/Microsoft.Extensions.Configuration.EnvironmentVariables/tests/EnvironmentVariablesTest.cs
//
// They were ported as-is to ensure our behaviour is identical to ASP.NET.

func setEnv(m map[string]string) map[string]string {
	if m == nil {
		m = make(map[string]string)
	}

	previous := make(map[string]string)
	for _, e := range os.Environ() {
		key := strings.Split(e, "=")[0]
		previous[key] = os.Getenv(key)
	}

	os.Clearenv()
	for k, v := range m {
		// Ignore errors... Even if this didn't work, it still keeps env vars in memory
		// so the subsequent os.Getenv(key) works just fine.
		//panic(err)
		_ = os.Setenv(k, v)
	}

	keys := os.Environ()
	fmt.Printf("%v", keys)

	return previous
}

func Test_envVarsConfigProvider_Load_LoadKeyValuePairsFromEnvironmentDictionary(t *testing.T) {
	in := map[string]string{
		"DefaultConnection:ConnectionString": "TestConnectionString",
		"DefaultConnection:Provider":         "SqlClient",
		"Inventory:ConnectionString":         "AnotherTestConnectionString",
		"Inventory:Provider":                 "MySql",
	}
	prev := setEnv(in)
	defer setEnv(prev)

	//config, err := NewEnvVarsSource("").Build()
	config, err := NewEnvVarsMapSource("", in).Build()
	assert.NoError(t, err)

	assert.Equal(t, "TestConnectionString", config.Get("defaultconnection:ConnectionString"))
	assert.Equal(t, "SqlClient", config.Get("DEFAULTCONNECTION:PROVIDER"))
	assert.Equal(t, "AnotherTestConnectionString", config.Get("Inventory:CONNECTIONSTRING"))
	assert.Equal(t, "MySql", config.Get("Inventory:Provider"))
	//assert.Equal(t, "EnvironmentVariablesConfigurationProvider Prefix: ''", config.Name())
}

func Test_envVarsConfigProvider_Load_LoadKeyValuePairsFromEnvironmentDictionaryWithPrefix(t *testing.T) {
	in := map[string]string{
		"DefaultConnection:ConnectionString": "TestConnectionString",
		"DefaultConnection:Provider":         "SqlClient",
		"Inventory:ConnectionString":         "AnotherTestConnectionString",
		"Inventory:Provider":                 "MySql",
	}
	prev := setEnv(in)
	defer setEnv(prev)

	config, err := NewEnvVarsSource("DefaultConnection").Build()
	assert.NoError(t, err)

	assert.Equal(t, "TestConnectionString", config.Get("ConnectionString"))
	assert.Equal(t, "SqlClient", config.Get("Provider"))
	//assert.Equal(t, "EnvironmentVariablesConfigurationProvider Prefix: 'DefaultConnection'", config.Name())
}

func Test_envVarsConfigProvider_Load_LoadKeyValuePairsFromAzureEnvironment(t *testing.T) {
	in := map[string]string{
		"APPSETTING_AppName":  "TestAppName",
		"CUSTOMCONNSTR_db1":   "CustomConnStr",
		"SQLCONNSTR_db2":      "SQLConnStr",
		"MYSQLCONNSTR_db3":    "MySQLConnStr",
		"SQLAZURECONNSTR_db4": "SQLAzureConnStr",
		"CommonEnv":           "CommonEnvValue",
	}
	prev := setEnv(in)
	defer setEnv(prev)

	config, err := NewEnvVarsSource("").Build()
	assert.NoError(t, err)

	var value string

	assert.Equal(t, "TestAppName", config.Get("APPSETTING_AppName"))
	assert.False(t, config.TryGet("AppName", &value))

	assert.True(t, config.TryGet("APPSETTING_AppName", &value))
	assert.Equal(t, "TestAppName", value)

	assert.Equal(t, "CustomConnStr", config.Get("ConnectionStrings:db1"))
	assert.Equal(t, "SQLConnStr", config.Get("ConnectionStrings:db2"))
	assert.Equal(t, "System.Data.SqlClient", config.Get("ConnectionStrings:db2_ProviderName"))
	assert.Equal(t, "MySQLConnStr", config.Get("ConnectionStrings:db3"))
	assert.Equal(t, "MySql.Data.MySqlClient", config.Get("ConnectionStrings:db3_ProviderName"))
	assert.Equal(t, "SQLAzureConnStr", config.Get("ConnectionStrings:db4"))
	assert.Equal(t, "System.Data.SqlClient", config.Get("ConnectionStrings:db4_ProviderName"))
	assert.Equal(t, "CommonEnvValue", config.Get("CommonEnv"))
}

func Test_envVarsConfigProvider_Load_LoadKeyValuePairsFromAzureEnvironmentWithPrefix(t *testing.T) {
	in := map[string]string{
		"CUSTOMCONNSTR_db1":   "CustomConnStr",
		"SQLCONNSTR_db2":      "SQLConnStr",
		"MYSQLCONNSTR_db3":    "MySQLConnStr",
		"SQLAZURECONNSTR_db4": "SQLAzureConnStr",
		"CommonEnv":           "CommonEnvValue",
	}
	prev := setEnv(in)
	defer setEnv(prev)

	config, err := NewEnvVarsSource("ConnectionStrings:").Build()
	assert.NoError(t, err)

	assert.Equal(t, "CustomConnStr", config.Get("db1"))
	assert.Equal(t, "SQLConnStr", config.Get("db2"))
	assert.Equal(t, "System.Data.SqlClient", config.Get("db2_ProviderName"))
	assert.Equal(t, "MySql.Data.MySqlClient", config.Get("db3_ProviderName"))
	assert.Equal(t, "System.Data.SqlClient", config.Get("db4_ProviderName"))
}

func Test_envVarsConfigProvider_Load_LastVariableAddedWhenKeyIsDuplicatedInAzureEnvironment(t *testing.T) {
	in := map[string]string{
		"ConnectionStrings:db2": "CommonEnvValue",
		"SQLCONNSTR_db2":        "SQLConnStr",
	}
	prev := setEnv(in)
	defer setEnv(prev)

	config, err := NewEnvVarsSource("").Build()
	assert.NoError(t, err)

	assert.NotEmpty(t, config.Get("ConnectionStrings:db2"))
	assert.Equal(t, "System.Data.SqlClient", config.Get("ConnectionStrings:db2_ProviderName"))
}

func Test_envVarsConfigProvider_Load_LastVariableAddedWhenMultipleEnvironmentVariablesWithSameNameButDifferentCaseExist(t *testing.T) {
	in := map[string]string{
		"CommonEnv": "CommonEnvValue1",
		"commonenv": "commonenvValue2",
		"cOMMonEnv": "commonenvValue3",
	}
	prev := setEnv(in)
	defer setEnv(prev)

	config, err := NewEnvVarsSource("").Build()
	assert.NoError(t, err)

	assert.NotEmpty(t, config.Get("cOMMonEnv"))
	assert.NotEmpty(t, config.Get("CommonEnv"))

}

func Test_envVarsConfigProvider_Load_ReplaceDoubleUnderscoreInEnvironmentVariables(t *testing.T) {
	in := map[string]string{
		"data__ConnectionString": "connection",
		"SQLCONNSTR_db1":         "connStr",
	}
	prev := setEnv(in)
	defer setEnv(prev)

	config, err := NewEnvVarsSource("").Build()
	assert.NoError(t, err)

	assert.Equal(t, "connection", config.Get("data:ConnectionString"))
	assert.Equal(t, "System.Data.SqlClient", config.Get("ConnectionStrings:db1_ProviderName"))
}

func Test_envVarsConfigProvider_Load_ReplaceDoubleUnderscoreInEnvironmentVariablesButNotPrefix(t *testing.T) {
	// Here we should get:
	//	 prefix:         "test__prefix__with__double__underscores__"
	//	 normalized key: "test:prefix:with:double:underscores:data:connectionstring
	//
	// Because "__" are not replaced in the prefix, the env var does is not
	// included in the final set of values.
	//
	in := map[string]string{
		"test__prefix__with__double__underscores__data__ConnectionString": "connection",
	}
	prev := setEnv(in)
	defer setEnv(prev)

	config, err := NewEnvVarsSource("test__prefix__with__double__underscores__").Build()
	assert.NoError(t, err)

	assert.Equal(t, "", config.Get("data:ConnectionString"))
}

func Test_envVarsConfigProvider_Load_ReplaceDoubleUnderscoreInEnvironmentVariablesButNotInAnomalousPrefix(t *testing.T) {
	in := map[string]string{
		"_____EXPERIMENTAL__data__ConnectionString": "connection",
	}
	prev := setEnv(in)
	defer setEnv(prev)

	config, err := NewEnvVarsSource("::_EXPERIMENTAL:").Build()
	assert.NoError(t, err)

	assert.Equal(t, "connection", config.Get("data:ConnectionString"))
}

func Test_envVarsConfigProvider_Load_ReplaceDoubleUnderscoreInEnvironmentVariablesWithDuplicatedPrefix(t *testing.T) {
	in := map[string]string{
		"test__test__ConnectionString": "connection",
	}
	prev := setEnv(in)
	defer setEnv(prev)

	config, err := NewEnvVarsSource("test__").Build()
	assert.NoError(t, err)

	// Original .NET code expects exception
	// Assert.Throws<InvalidOperationException>(() => config.Get("test:ConnectionString"));
	assert.Equal(t, "", config.Get("test:ConnectionString"))
}

func Test_envVarsConfigProvider_Load_PrefixPreventsLoadingSqlConnectionStrings(t *testing.T) {
	in := map[string]string{
		"test__test__ConnectionString": "connection",
		"SQLCONNSTR_db1":               "connStr",
	}
	prev := setEnv(in)
	defer setEnv(prev)

	config, err := NewEnvVarsSource("test:").Build()
	assert.NoError(t, err)

	assert.Equal(t, "connection", config.Get("test:ConnectionString"))

	// Original .NET code expects exception
	// Assert.Throws<InvalidOperationException>(() => config.Get("ConnectionStrings:db1_ProviderName"));
	assert.Equal(t, "", config.Get("ConnectionStrings:db1_ProviderName"))
}

// ------- END of CORE TESTS -------

// TODO: Add additional tests when we have interface ready.
//	- AddEnvironmentVariables_Bind_PrefixShouldNormalize
//	- AddEnvironmentVariables_UsingDoubleUnderscores_Bind_PrefixWontNormalize
//	- BindingDoesNotThrowIfReloadedDuringBinding
//

// ------- Extra custom tests  -------

func Test_envVarsConfigProvider_Load_LoadsFromEnv(t *testing.T) {
	prev := setEnv(nil)
	defer setEnv(prev)

	var err error

	os.Clearenv()

	err = os.Setenv("myFoo", "foo value = with equal sign")
	assert.NoError(t, err)

	err = os.Setenv("myBar", "bar value = with equal sign")
	assert.NoError(t, err)

	config, err := NewEnvVarsSource("").Build()
	assert.NoError(t, err)

	assert.Equal(t, "foo value = with equal sign", config.Get("myFoo"))
	assert.Equal(t, "bar value = with equal sign", config.Get("myBar"))

	// Load reloads everything
	err = os.Setenv("myBar", "changed value")
	assert.NoError(t, err)

	config, err = NewEnvVarsSource("").Build()
	assert.NoError(t, err)

	assert.Equal(t, "foo value = with equal sign", config.Get("myFoo"))
	assert.Equal(t, "changed value", config.Get("myBar"))

}
