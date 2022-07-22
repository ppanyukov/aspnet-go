// Package config implements ASP.NET equivalents of Microsoft.Extensions.Configuration.*
// libraries which provide configuration from various sources such as environment
// variables, json files etc.
//
// See: https://github.com/dotnet/runtime/tree/release/6.0/src/libraries/Microsoft.Extensions.Configuration/src
//
// Standard order of config in ASP.NET, the sources added later override
// settings in previous sources:
//
//		appConfigBuilder.AddJsonFile("appsettings.json", optional: true, reloadOnChange: reloadOnChange)
//		appConfigBuilder.AddJsonFile($"appsettings.{env.EnvironmentName}.json", optional: true, reloadOnChange: reloadOnChange);
//		if (env.IsDevelopment()) {
//			appConfigBuilder.AddUserSecrets(appAssembly, optional: true, reloadOnChange: reloadOnChange);
//		}
//		appConfigBuilder.AddEnvironmentVariables();
//		appConfigBuilder.AddCommandLine(args);
//
// See: https://github.com/dotnet/runtime/blob/backport/pr-69503-to-release/6.0/src/libraries/Microsoft.Extensions.Hosting/src/HostingHostBuilderExtensions.cs#L202
package config

import "strings"

// keyDelimiter is a hierarchical delimiter for keys.
const keyDelimiter = ":"

func normalizeKey(k string) string {
	k = strings.ToLower(k)
	k = strings.Replace(k, "__", keyDelimiter, -1)
	return k
}
