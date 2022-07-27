package config_test

import (
	"fmt"
	"log"

	"github.com/ppanyukov/aspnet-go/pkg/config"
)

func Example_envVarsPrefix() {
	// In ASP.NET Environment Variables provider you can specify a prefix.
	// In such case only env vars with specified prefix are loaded. The prefix
	// is also stripped away from the key. This behaviour is preserved here.
	const prefix = "MYAPP_"

	env := map[string]string{
		// These will be included and the prefix will be stripped from key names.
		"MYAPP_SETTING_A": "A",
		"MYAPP_SETTING_B": "B",

		// These will be excluded
		"OTHERAPP_SETTING_A": "A",
		"OTHERAPP_SETTING_B": "A",
	}

	builder := config.NewBuilder()
	builder.AddSource(config.NewEnvVarsMapSource(prefix, env))
	c, err := builder.Build()
	if err != nil {
		log.Fatal(err)
	}

	for _, key := range c.Keys() {
		val := c.Get(key)
		fmt.Printf("key=%s, value=%s\n", key, val)
	}

	// Output:
	// key=setting_a, value=A
	// key=setting_b, value=B
}
