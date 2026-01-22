//go:build gin

package simple

import (
	"net/http"

	ginadapter "github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
	ginlib "github.com/gin-gonic/gin"
)

// GinRouter wraps the gin adapter Router and injects options from Spec automatically.
type GinRouter struct {
	Base *ginadapter.Router
	Spec Spec
}

func NewGin(base *ginadapter.Router, spec Spec) *GinRouter {
	return &GinRouter{Base: base, Spec: spec}
}

func (r *GinRouter) Routes() []openapi.RouteMeta { return r.Base.Routes() }

func (r *GinRouter) Group(prefix string, opts ...ginadapter.HandlerOption) *ginadapter.Group {
	return r.Base.Group(prefix, opts...)
}

func (r *GinRouter) Handle(method, path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	all := make([]openapi.HandlerOption, 0, len(opts))
	for _, o := range opts {
		all = append(all, o)
	}
	if def, ok := r.Spec[Key(method, path)]; ok {
		all = Inject(all, def)
	}
	// Convert back to adapter options (same underlying type)
	out := make([]ginadapter.HandlerOption, 0, len(all))
	for _, o := range all {
		out = append(out, o)
	}
	r.Base.Handle(method, path, h, out...)
}

func (r *GinRouter) GET(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	r.Handle(http.MethodGet, path, h, opts...)
}
func (r *GinRouter) POST(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	r.Handle(http.MethodPost, path, h, opts...)
}
func (r *GinRouter) PUT(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	r.Handle(http.MethodPut, path, h, opts...)
}
func (r *GinRouter) PATCH(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	r.Handle(http.MethodPatch, path, h, opts...)
}
func (r *GinRouter) DELETE(path string, h ginlib.HandlerFunc, opts ...ginadapter.HandlerOption) {
	r.Handle(http.MethodDelete, path, h, opts...)
}
