// Package profile manages named environment profiles (e.g. dev, staging, prod)
// that group layers and validation schemas together for a given deployment target.
package profile

import (
	"errors"
	"fmt"
)

// Profile represents a named deployment target with associated layer names
// and required validation schema keys.
type Profile struct {
	Name       string
	Layers     []string
	RequiredKeys []string
	tags       map[string]string
}

// NewProfile creates a new Profile with the given name.
// Returns an error if the name is empty.
func NewProfile(name string) (*Profile, error) {
	if name == "" {
		return nil, errors.New("profile name must not be empty")
	}
	return &Profile{
		Name: name,
		tags: make(map[string]string),
	}, nil
}

// AddLayer appends a layer name to the profile's ordered layer list.
// Duplicate layer names are silently ignored.
func (p *Profile) AddLayer(layer string) error {
	if layer == "" {
		return errors.New("layer name must not be empty")
	}
	for _, l := range p.Layers {
		if l == layer {
			return nil
		}
	}
	p.Layers = append(p.Layers, layer)
	return nil
}

// RequireKey marks a key as required for this profile during validation.
func (p *Profile) RequireKey(key string) error {
	if key == "" {
		return errors.New("required key must not be empty")
	}
	for _, k := range p.RequiredKeys {
		if k == key {
			return nil
		}
	}
	p.RequiredKeys = append(p.RequiredKeys, key)
	return nil
}

// SetTag attaches an arbitrary metadata tag to the profile.
func (p *Profile) SetTag(key, value string) error {
	if key == "" {
		return errors.New("tag key must not be empty")
	}
	p.tags[key] = value
	return nil
}

// Tag returns the value of a metadata tag and whether it was found.
func (p *Profile) Tag(key string) (string, bool) {
	v, ok := p.tags[key]
	return v, ok
}

// Validate checks that the profile is internally consistent.
func (p *Profile) Validate() error {
	if len(p.Layers) == 0 {
		return fmt.Errorf("profile %q has no layers defined", p.Name)
	}
	return nil
}
