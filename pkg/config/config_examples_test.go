package config_test

import (
	"fmt"
	"github.com/ppanyukov/aspnet-go/pkg/config"
)

func ExampleBuilder_Build_jsonMerge() {
	// Let's say have these two Json config files: appsettings.json and appsettings.Development.json.
	//
	//The appsettings.json contains common settings, and appsettings.Development.json contains
	// settings specific for the Development environment. Any settings in
	// appsettings.Development.json should override settings in appsettings.json.
	//
	// Note that comments are allowed as they are allowed in ASP.NET too.
	jsonFileAppsettings := `
		{
			// AppName is the same between environments
			"AppName": "app name from the appsettings.json",

			// 
			"ConnectionsStrings": {
				"Redis": "redis connection string from appsettings.json",
				"Sql": "sql connection string from appsettings.json"
			}
		}
	`

	jsonFileAppsettingsDevelopment := `
		{
			"ConnectionsStrings": {
				// Override 
				"Redis": "redis connection string from appsettings.Development.json",

				// Extra settings to use this logger in development environment
				"Logger": "logger connection string from appsettings.Development.json"
			}
		}
	`

	// Builder allows us to build a configuration from multiple sources.
	// The later the source added the higher its precedence.
	//
	// When built, the configuration will be equivalent to the following JSON file.
	_ = `
		{
			"AppName": "app name from the appsettings.json",
			"ConnectionsStrings": {
				"Redis": "redis connection string from appsettings.Development.json",
				"Sql": "sql connection string from appsettings.json",
				"Logger": "logger connection string from appsettings.Development.json"
			}
		}
	`
	builder := config.NewBuilder()
	builder.AddSource(config.NewJsonSource([]byte(jsonFileAppsettings)))
	builder.AddSource(config.NewJsonSource([]byte(jsonFileAppsettingsDevelopment)))

	rootConfig, err := builder.Build()
	if err != nil {
		panic(err)
	}

	// Print out key value pairs in the merged config.
	// The key names are normalised to lower case, and ":" separator is used
	// to denote key structure in the JSON file.
	for _, key := range rootConfig.Keys() {
		val := rootConfig.Get(key)
		fmt.Printf("%s = %s\n", key, val)
	}

	// Output:
	// appname = app name from the appsettings.json
	// connectionsstrings:logger = logger connection string from appsettings.Development.json
	// connectionsstrings:redis = redis connection string from appsettings.Development.json
	// connectionsstrings:sql = sql connection string from appsettings.json
}
