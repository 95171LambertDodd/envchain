// Package pin implements key pinning for envchain.
//
// A Pinner records the expected value for one or more environment keys and
// validates them against any Source at runtime. This is useful for detecting
// configuration drift between environments: pin the authoritative baseline
// (e.g. the production layer) and then validate the composed chain to ensure
// no higher-priority layer has silently overridden a critical key.
//
// Basic usage:
//
//	p := pin.NewPinner()
//	_ = p.Pin("APP_ENV", "production")
//	if errs := p.Validate(src); len(errs) > 0 {
//		for _, e := range errs { log.Println(e) }
//	}
//
// Layer helpers:
//
//	p, err := pin.PinFromLayer(baseLayer)
//	errs := pin.ValidateChain(p, chain)
package pin
