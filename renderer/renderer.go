package renderer

import (
	"git.area51.dev/peter/videoident/util"
)

// Renderer that can be registered
type Renderer func(Context) error

// Of returns a Renderer that will call each Renderer provided.
func Of(rr ...Renderer) Renderer {
	switch len(rr) {
	case 0:
		return nil
	case 1:
		return rr[0]
	default:
		var r Renderer
		for _, e := range rr {
			r = r.Then(e)
		}
		return r
	}
}

// Do calls a renderer whilst being nil safe.
// This is the same as calling r(ctx) however if r==nil then this
// will do nothing whilst r(ctx) will panic
func (r Renderer) Do(ctx Context) error {
	if r != nil {
		return r(ctx)
	}
	return nil
}

// Then returns a Renderer that calls this renderer then a second one.
func (r Renderer) Then(b Renderer) Renderer {
	if r == nil {
		return b
	}
	if b == nil {
		return r
	}
	return func(ctx Context) error {
		if err := r(ctx); err != nil {
			return err
		}
		return b(ctx)
	}
}

// Within ensures that the renderer will only be invoked
// for frames within two boundaries
func (r Renderer) Within(s, e int) Renderer {
	if r == nil {
		return nil
	}

	return func(ctx Context) error {
		if util.Within(ctx.Frame(), s, e) {
			return r.Do(ctx)
		}
		return nil
	}
}
