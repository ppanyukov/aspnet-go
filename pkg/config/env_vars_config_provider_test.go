package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// These tests are a port of EnvironmentVariablesTest.cs
// https://github.com/dotnet/runtime/blob/release/6.0/src/libraries/Microsoft.Extensions.Configuration.EnvironmentVariables/tests/EnvironmentVariablesTest.cs
//
func Test_envVarsConfigProvider_Load_LoadKeyValuePairsFromEnvironmentDictionary(t *testing.T) {
	in := map[string]string{
		"DefaultConnection:ConnectionString": "TestConnectionString",
		"DefaultConnection:Provider":         "SqlClient",
		"Inventory:ConnectionString":         "AnotherTestConnectionString",
		"Inventory:Provider":                 "MySql",
	}

	envConfigSrc := newEnvVars("")
	envConfigSrc.load(in)

	assert.Equal(t, "TestConnectionString", envConfigSrc.Get("defaultconnection:ConnectionString"))
	assert.Equal(t, "SqlClient", envConfigSrc.Get("DEFAULTCONNECTION:PROVIDER"))
	assert.Equal(t, "AnotherTestConnectionString", envConfigSrc.Get("Inventory:CONNECTIONSTRING"))
	assert.Equal(t, "MySql", envConfigSrc.Get("Inventory:Provider"))
	assert.Equal(t, "EnvironmentVariablesConfigurationProvider Prefix: ''", envConfigSrc.String())
}

func Test_envVarsConfigProvider_Load_LoadKeyValuePairsFromEnvironmentDictionaryWithPrefix(t *testing.T) {
	in := map[string]string{
		"DefaultConnection:ConnectionString": "TestConnectionString",
		"DefaultConnection:Provider":         "SqlClient",
		"Inventory:ConnectionString":         "AnotherTestConnectionString",
		"Inventory:Provider":                 "MySql",
	}

	envConfigSrc := newEnvVars("DefaultConnection")
	envConfigSrc.load(in)

	assert.Equal(t, "TestConnectionString", envConfigSrc.Get("ConnectionString"))
	assert.Equal(t, "SqlClient", envConfigSrc.Get("Provider"))
	assert.Equal(t, "EnvironmentVariablesConfigurationProvider Prefix: 'DefaultConnection'", envConfigSrc.String())
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

	envConfigSrc := newEnvVars("")
	envConfigSrc.load(in)

	var value string

	assert.Equal(t, "TestAppName", envConfigSrc.Get("APPSETTING_AppName"))
	assert.False(t, envConfigSrc.TryGet("AppName", &value))

	assert.True(t, envConfigSrc.TryGet("APPSETTING_AppName", &value))
	assert.Equal(t, "TestAppName", value)

	assert.Equal(t, "CustomConnStr", envConfigSrc.Get("ConnectionStrings:db1"))
	assert.Equal(t, "SQLConnStr", envConfigSrc.Get("ConnectionStrings:db2"))
	assert.Equal(t, "System.Data.SqlClient", envConfigSrc.Get("ConnectionStrings:db2_ProviderName"))
	assert.Equal(t, "MySQLConnStr", envConfigSrc.Get("ConnectionStrings:db3"))
	assert.Equal(t, "MySql.Data.MySqlClient", envConfigSrc.Get("ConnectionStrings:db3_ProviderName"))
	assert.Equal(t, "SQLAzureConnStr", envConfigSrc.Get("ConnectionStrings:db4"))
	assert.Equal(t, "System.Data.SqlClient", envConfigSrc.Get("ConnectionStrings:db4_ProviderName"))
	assert.Equal(t, "CommonEnvValue", envConfigSrc.Get("CommonEnv"))
}

func Test_envVarsConfigProvider_Load_LoadKeyValuePairsFromAzureEnvironmentWithPrefix(t *testing.T) {
	in := map[string]string{
		"CUSTOMCONNSTR_db1":   "CustomConnStr",
		"SQLCONNSTR_db2":      "SQLConnStr",
		"MYSQLCONNSTR_db3":    "MySQLConnStr",
		"SQLAZURECONNSTR_db4": "SQLAzureConnStr",
		"CommonEnv":           "CommonEnvValue",
	}

	envConfigSrc := newEnvVars("ConnectionStrings:")
	envConfigSrc.load(in)

	assert.Equal(t, "CustomConnStr", envConfigSrc.Get("db1"))
	assert.Equal(t, "SQLConnStr", envConfigSrc.Get("db2"))
	assert.Equal(t, "System.Data.SqlClient", envConfigSrc.Get("db2_ProviderName"))
	assert.Equal(t, "MySql.Data.MySqlClient", envConfigSrc.Get("db3_ProviderName"))
	assert.Equal(t, "System.Data.SqlClient", envConfigSrc.Get("db4_ProviderName"))
}

func Test_envVarsConfigProvider_Load_LastVariableAddedWhenKeyIsDuplicatedInAzureEnvironment(t *testing.T) {
	in := map[string]string{
		"ConnectionStrings:db2": "CommonEnvValue",
		"SQLCONNSTR_db2":        "SQLConnStr",
	}

	envConfigSrc := newEnvVars("")
	envConfigSrc.load(in)

	assert.NotEmpty(t, envConfigSrc.Get("ConnectionStrings:db2"))
	assert.Equal(t, "System.Data.SqlClient", envConfigSrc.Get("ConnectionStrings:db2_ProviderName"))
}

func Test_envVarsConfigProvider_Load_LastVariableAddedWhenMultipleEnvironmentVariablesWithSameNameButDifferentCaseExist(t *testing.T) {
	in := map[string]string{
		"CommonEnv": "CommonEnvValue1",
		"commonenv": "commonenvValue2",
		"cOMMonEnv": "commonenvValue3",
	}

	envConfigSrc := newEnvVars("")
	envConfigSrc.load(in)

	assert.NotEmpty(t, envConfigSrc.Get("cOMMonEnv"))
	assert.NotEmpty(t, envConfigSrc.Get("CommonEnv"))

}

func Test_envVarsConfigProvider_Load_ReplaceDoubleUnderscoreInEnvironmentVariables(t *testing.T) {
	in := map[string]string{
		"data__ConnectionString": "connection",
		"SQLCONNSTR_db1":         "connStr",
	}

	envConfigSrc := newEnvVars("")
	envConfigSrc.load(in)

	assert.Equal(t, "connection", envConfigSrc.Get("data:ConnectionString"))
	assert.Equal(t, "System.Data.SqlClient", envConfigSrc.Get("ConnectionStrings:db1_ProviderName"))
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

	envConfigSrc := newEnvVars("test__prefix__with__double__underscores__")
	envConfigSrc.load(in)

	assert.Equal(t, "", envConfigSrc.Get("data:ConnectionString"))
}

func Test_envVarsConfigProvider_Load_ReplaceDoubleUnderscoreInEnvironmentVariablesButNotInAnomalousPrefix(t *testing.T) {
	in := map[string]string{
		"_____EXPERIMENTAL__data__ConnectionString": "connection",
	}

	envConfigSrc := newEnvVars("::_EXPERIMENTAL:")
	envConfigSrc.load(in)

	assert.Equal(t, "connection", envConfigSrc.Get("data:ConnectionString"))
}

func Test_envVarsConfigProvider_Load_ReplaceDoubleUnderscoreInEnvironmentVariablesWithDuplicatedPrefix(t *testing.T) {
	in := map[string]string{
		"test__test__ConnectionString": "connection",
	}

	envConfigSrc := newEnvVars("test__")
	envConfigSrc.load(in)

	// Original .NET code expects exception
	// Assert.Throws<InvalidOperationException>(() => envConfigSrc.Get("test:ConnectionString"));
	assert.Equal(t, "", envConfigSrc.Get("test:ConnectionString"))
}

func Test_envVarsConfigProvider_Load_PrefixPreventsLoadingSqlConnectionStrings(t *testing.T) {
	in := map[string]string{
		"test__test__ConnectionString": "connection",
		"SQLCONNSTR_db1":               "connStr",
	}

	envConfigSrc := newEnvVars("test:")
	envConfigSrc.load(in)

	assert.Equal(t, "connection", envConfigSrc.Get("test:ConnectionString"))

	// Original .NET code expects exception
	// Assert.Throws<InvalidOperationException>(() => envConfigSrc.Get("ConnectionStrings:db1_ProviderName"));
	assert.Equal(t, "", envConfigSrc.Get("ConnectionStrings:db1_ProviderName"))
}

// ------- END of CORE TESTS -------

// TODO: Add additional tests when we have interface ready.
//	- AddEnvironmentVariables_Bind_PrefixShouldNormalize
//	- AddEnvironmentVariables_UsingDoubleUnderscores_Bind_PrefixWontNormalize
//	- BindingDoesNotThrowIfReloadedDuringBinding
//

// ------- Extra custom tests  -------

func Test_envVarsConfigProvider_Load_LoadsFromEnv(t *testing.T) {
	var err error

	err = os.Setenv("myFoo", "foo value = with equal sign")
	assert.NoError(t, err)

	err = os.Setenv("myBar", "bar value = with equal sign")
	assert.NoError(t, err)

	envConfigSrc := newEnvVars("")

	envConfigSrc.Load()
	assert.Equal(t, "foo value = with equal sign", envConfigSrc.Get("myFoo"))
	assert.Equal(t, "bar value = with equal sign", envConfigSrc.Get("myBar"))

	// Load reloads everything
	err = os.Setenv("myBar", "changed value")
	assert.NoError(t, err)
	envConfigSrc.Load()
	assert.Equal(t, "foo value = with equal sign", envConfigSrc.Get("myFoo"))
	assert.Equal(t, "changed value", envConfigSrc.Get("myBar"))

}
