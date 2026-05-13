// Package env bridges the OS environment with envchain's layered config model.
//
// # Overview
//
// env.Loader reads key=value pairs from os.Environ and populates a
// config.Layer. Two options control its behaviour:
//
//   - WithPrefix(p) — only variables whose names start with p are loaded;
//     the prefix is stripped before storing so that APP_HOST becomes HOST.
//
//   - WithStrict(keys...) — after loading, each listed key must be present
//     in the resulting layer or Load returns an error. Keys are matched
//     after prefix stripping.
//
// # Example
//
//	loader := env.NewLoader(
//		env.WithPrefix("APP_"),
//		env.WithStrict("HOST", "PORT"),
//	)
//	layer, err := loader.Load("app-env")
//	if err != nil {
//		log.Fatal(err)
//	}
package env
