package interpolate

// ChainResolver returns a Resolver that queries each provided Resolver in
// order, returning the first match found.
func ChainResolver(resolvers ...Resolver) Resolver {
	return func(key string) (string, bool) {
		for _, r := range resolvers {
			if v, ok := r(key); ok {
				return v, true
			}
		}
		return "", false
	}
}

// MapResolver builds a Resolver backed by a plain string map.
func MapResolver(m map[string]string) Resolver {
	return func(key string) (string, bool) {
		v, ok := m[key]
		return v, ok
	}
}
