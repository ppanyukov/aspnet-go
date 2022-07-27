package config_test

import (
	"fmt"
	"log"

	"github.com/ppanyukov/aspnet-go/pkg/config"
)

func Example_keyHierarchy() {
	// ASP.NET config supports hierarchical keys, especially when it comes to Json.
	// Env vars also have this support using a special double underscore "__"
	// as a separator.

	// Consider a Json config file like this.
	json := `
	{
		// Since ASP.NET config allows comments in JSON file so do we.
		"MyApp": {
			"Logging": {
				// This will be under "myapp:logging:enabled" key.
				// The value will be a string, not a bool.
				"Enabled": true,

				// This will be under "myapp:logging:level" key
				"Level": "info"
			}
		}
	}
	`

	// We can override Json config with env vars like so.
	// Note that the casing of keys doesn't matter.
	env := map[string]string{
		// This also will be under "myapp:logging:level" key
		// and will override the setting in json file.
		"MYAPP__LOGGING__level": "debug",
	}

	builder := config.NewBuilder()
	builder.AddSource(config.NewJsonSource([]byte(json)).WithName("json file"))
	builder.AddSource(config.NewEnvVarsMapSource("", env).WithName("env vars"))
	c, err := builder.Build()
	if err != nil {
		log.Fatal(err)
	}

	for _, key := range c.Keys() {
		entry := c.GetEntry(key)
		fmt.Printf("key=%s, value=%s, source=%s\n", key, entry.Value(), entry.Source().Name())
	}

	// Output:
	// key=myapp:logging:enabled, value=true, source=json file
	// key=myapp:logging:level, value=debug, source=env vars
}
