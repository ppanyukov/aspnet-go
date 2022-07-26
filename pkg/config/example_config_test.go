package config_test

import (
	"fmt"
	"log"

	"github.com/ppanyukov/aspnet-go/pkg/config"
)

var json = []byte(`{
		"foo": "foo from json",
		"bar": "bar from json",
		"ConnectionStrings": {
			"Sql": "sql from json",
			"Redis": "redis from json"
		}
	}`)

var env = map[string]string{
	"foo":                     "foo from env",
	"zoo":                     "zoo from env",
	"CONNECTIONSTRINGS__SQL":  "sql from env",
	"CONNECTIONSTRINGS__BLOB": "blob from env",
}

func Example() {
	// TODO(ppanyukov): write more comments and more examples.
	builder := config.NewBuilder()
	builder.AddSource(config.NewJsonSource([]byte(json)).WithName("json provider"))
	builder.AddSource(config.NewEnvVarsMapSource("", env).WithName("env provider"))
	c, err := builder.Build()
	if err != nil {
		log.Fatal(err)
	}

	// We can get the list of all keys in settings.
	fmt.Printf("KEYS - sorted alphabetically:\n")
	for _, key := range c.Keys() {
		fmt.Printf("  - %s\n", key)
	}

	// List of entries allows to iterate through all settings
	fmt.Printf("ENTRIES - sorted by key:\n")
	for _, entry := range c.GetEntries() {
		fmt.Printf("  - key=%s, value=%v, provider=%v\n", entry.Key(), entry.Value(), entry.Source().Name())
	}

	// We can get individual values for each key
	fmt.Printf("KEY VALUES:\n")
	for _, key := range c.Keys() {
		val := c.Get(key)
		fmt.Printf("  - key=%s, value=%v\n", key, val)
	}

	// key=foo, value=foo value from env, provider=env vars provider
	// key=bar, value=bar value from json, provider=json provider

	// Output:
	// KEYS - sorted alphabetically:
	//   - bar
	//   - connectionstrings:blob
	//   - connectionstrings:redis
	//   - connectionstrings:sql
	//   - foo
	//   - zoo
	// ENTRIES - sorted by key:
	//   - key=bar, value=bar from json, provider=json provider
	//   - key=connectionstrings:blob, value=blob from env, provider=env provider
	//   - key=connectionstrings:redis, value=redis from json, provider=json provider
	//   - key=connectionstrings:sql, value=sql from env, provider=env provider
	//   - key=foo, value=foo from env, provider=env provider
	//   - key=zoo, value=zoo from env, provider=env provider
	// KEY VALUES:
	//   - key=bar, value=bar from json
	//   - key=connectionstrings:blob, value=blob from env
	//   - key=connectionstrings:redis, value=redis from json
	//   - key=connectionstrings:sql, value=sql from env
	//   - key=foo, value=foo from env
	//   - key=zoo, value=zoo from env

}
