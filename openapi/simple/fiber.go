//go:build fiber

package simple

import (
	"net/http"

	fiberadapter "github.com/aizacoders/openapigo/adapters/fiber"
	"github.com/aizacoders/openapigo/openapi"
	fiberlib "github.com/gofiber/fiber/v2"
)

// FiberRouter wraps the fiber adapter Router and injects options from Spec automatically.
type FiberRouter struct {
	Base *fiberadapter.Router
	Spec Spec
}

func NewFiber(base *fiberadapter.Router, spec Spec) *FiberRouter {
	return &FiberRouter{Base: base, Spec: spec}
}

func (r *FiberRouter) Routes() []openapi.RouteMeta { return r.Base.Routes() }

func (r *FiberRouter) Group(prefix string, opts ...fiberadapter.HandlerOption) *fiberadapter.Group {
	return r.Base.Group(prefix, opts...)
}

func (r *FiberRouter) Handle(method, path string, h fiberlib.Handler, opts ...fiberadapter.HandlerOption) {
	all := make([]openapi.HandlerOption, 0, len(opts))
	for _, o := range opts {
		all = append(all, o)
	}
	if def, ok := r.Spec[Key(method, path)]; ok {
		all = Inject(all, def)
	}
	out := make([]fiberadapter.HandlerOption, 0, len(all))
	for _, o := range all {
		out = append(out, o)
	}
	r.Base.Handle(method, path, h, out...)
}

func (r *FiberRouter) GET(path string, h fiberlib.Handler, opts ...fiberadapter.HandlerOption) {
	r.Handle(http.MethodGet, path, h, opts...)
}
func (r *FiberRouter) POST(path string, h fiberlib.Handler, opts ...fiberadapter.HandlerOption) {
	r.Handle(http.MethodPost, path, h, opts...)
}
func (r *FiberRouter) PUT(path string, h fiberlib.Handler, opts ...fiberadapter.HandlerOption) {
	r.Handle(http.MethodPut, path, h, opts...)
}
func (r *FiberRouter) PATCH(path string, h fiberlib.Handler, opts ...fiberadapter.HandlerOption) {
	r.Handle(http.MethodPatch, path, h, opts...)
}
func (r *FiberRouter) DELETE(path string, h fiberlib.Handler, opts ...fiberadapter.HandlerOption) {
	r.Handle(http.MethodDelete, path, h, opts...)
}
