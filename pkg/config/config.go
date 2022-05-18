// Package config implements ASP.NET equivalents of Microsoft.Extensions.Configuration.*
// libraries which provide configuration from various sources such as environment
// variables, json files etc.
// See: https://github.com/dotnet/runtime/tree/release/6.0/src/libraries/Microsoft.Extensions.Configuration/src
package config

import "strings"

// keyDelimiter is a hierarchical delimiter for keys.
const keyDelimiter = ":"

func normalizeKey(k string) string {
	k = strings.ToLower(k)
	k = strings.Replace(k, "__", keyDelimiter, -1)
	return k
}
