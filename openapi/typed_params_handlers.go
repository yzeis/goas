package openapi

import (
	"net/http"
	"reflect"
	"strings"
)

// Param binds a name (path/query) to a typed scalar.
//
// You can use it in typed handler signatures to get automatic, documented params.
// Example:
//
//	func(w http.ResponseWriter, r *http.Request, id openapi.Param[int]) (...)
//
// The router will fill Name and Value, and you can use id.Value.
// If parsing fails, the handler wrapper returns 400.
type Param[T any] struct {
	Name  string
	Value T
}

// --- typed path-param handlers ---

type PathTypedHandler1[TReq any, TRes any, TP1 any] func(w http.ResponseWriter, r *http.Request, req TReq, p1 Param[TP1]) (TRes, int, error)

// GETP1T registers a GET handler with one typed path param referenced by name.
func GETP1T[TReq any, TRes any, TP1 any](router *Router, path string, p1Name string, h PathTypedHandler1[TReq, TRes, TP1], opts ...HandlerOption) {
	paramOpt := WithPathParam(p1Name, kindToParamType(reflect.TypeOf(*new(TP1)).Kind()), true, "")

	// Auto schema options
	reqOpt, resOpt := typedOptions[TReq, TRes]()
	base := make([]HandlerOption, 0, 3)
	if reqOpt != nil {
		base = append(base, reqOpt)
	}
	if resOpt != nil {
		base = append(base, resOpt)
	}
	base = append(base, paramOpt)
	base = mergeOpts(base, opts...)

	router.GET(path, func(w http.ResponseWriter, r *http.Request) {
		var reqVal TReq
		if !isZeroStructType[TReq]() {
			_ = Bind(r, &reqVal)
		}

		pv, err := PathParamValue[TP1](r, p1Name)
		if err != nil {
			JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		res, code, err := h(w, r, reqVal, Param[TP1]{Name: p1Name, Value: pv.Value})
		if err != nil {
			JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if code == 0 {
			code = http.StatusOK
		}
		if isZeroStructType[TRes]() {
			w.WriteHeader(code)
			return
		}
		JSON(w, code, res)
	}, base...)
}

// POSTP1T registers a POST handler with one typed path param referenced by name.
func POSTP1T[TReq any, TRes any, TP1 any](router *Router, path string, p1Name string, h PathTypedHandler1[TReq, TRes, TP1], opts ...HandlerOption) {
	paramOpt := WithPathParam(p1Name, kindToParamType(reflect.TypeOf(*new(TP1)).Kind()), true, "")

	reqOpt, resOpt := typedOptions[TReq, TRes]()
	base := make([]HandlerOption, 0, 3)
	if reqOpt != nil {
		base = append(base, reqOpt)
	}
	if resOpt != nil {
		base = append(base, resOpt)
	}
	base = append(base, paramOpt)
	base = mergeOpts(base, opts...)

	router.POST(path, func(w http.ResponseWriter, r *http.Request) {
		var reqVal TReq
		if !isZeroStructType[TReq]() {
			_ = Bind(r, &reqVal)
		}

		pv, err := PathParamValue[TP1](r, p1Name)
		if err != nil {
			JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		res, code, err := h(w, r, reqVal, Param[TP1]{Name: p1Name, Value: pv.Value})
		if err != nil {
			JSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if code == 0 {
			code = http.StatusOK
		}
		if isZeroStructType[TRes]() {
			w.WriteHeader(code)
			return
		}
		JSON(w, code, res)
	}, base...)
}

// --- query inference v2 (typed struct) ---

// Query is a typed specification for query parameters. It’s a small wrapper type
// so we can distinguish it from request body structs.
//
// Any exported field with a `query:"name"` tag becomes a query parameter.
// Supported field types: string, bool, ints, uints, floats, and slices of those.
// Required: `queryRequired:"true"` or `required:"true"`.
// Description: `desc:"..."`.
//
// Example:
//   type Q struct {
//     Q string `query:"q" queryRequired:"true"`
//     Limit int `query:"limit"`
//   }
//   openapi.WithQueryStruct(Q{})

type Query[T any] struct{ Value T }

func WithQueryStruct[T any](sample T) HandlerOption {
	params := inferQueryParamsFromStruct(sample)
	return WithQueryParams(params...)
}

func inferQueryParamsFromStruct(sample any) []QueryParam {
	t := reflect.TypeOf(sample)
	if t == nil {
		return nil
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	out := make([]QueryParam, 0)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		name := strings.TrimSpace(f.Tag.Get("query"))
		if name == "" || name == "-" {
			continue
		}
		req := strings.EqualFold(f.Tag.Get("queryRequired"), "true") || strings.EqualFold(f.Tag.Get("required"), "true")
		desc := f.Tag.Get("desc")

		ft := f.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		// slice support (treat as repeated query param)
		if ft.Kind() == reflect.Slice || ft.Kind() == reflect.Array {
			ft = ft.Elem()
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
		}

		out = append(out, QueryParam{
			Name:        name,
			Type:        kindToParamType(ft.Kind()),
			Required:    req,
			Description: desc,
		})
	}
	return out
}

func kindToParamType(k reflect.Kind) ParamType {
	switch k {
	case reflect.String:
		return ParamString
	case reflect.Bool:
		return ParamBoolean
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ParamInteger
	case reflect.Float32, reflect.Float64:
		return ParamNumber
	default:
		return ParamString
	}
}
