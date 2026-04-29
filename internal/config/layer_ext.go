package config

// All returns a copy of all key-value pairs stored in the layer.
// This is used by external packages (e.g. merge) to iterate over entries.
func (l *Layer) All() map[string]string {
	copy := make(map[string]string, len(l.data))
	for k, v := range l.data {
		copy[k] = v
	}
	return copy
}

// Name returns the layer's identifier.
func (l *Layer) Name() string {
	return l.name
}
