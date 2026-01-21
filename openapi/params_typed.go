package openapi

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// ParamType represents a primitive type used for path/query parameters.
type ParamType string

const (
	ParamString  ParamType = "string"
	ParamInteger ParamType = "integer"
	ParamNumber  ParamType = "number"
	ParamBoolean ParamType = "boolean"
)

// QueryParam describes a query parameter for OpenAPI generation.
// Name only (no regex). For advanced usage, we can add style/explode/etc later.
type QueryParam struct {
	Name        string
	Type        ParamType
	Required    bool
	Description string
}

// WithQueryParams declares query parameters for a route for OpenAPI generation.
//
// This is the first step for query inference. Purely "automatic" query inference
// from handler code isn't reliable in Go, so we offer a small declarative API.
func WithQueryParams(params ...QueryParam) HandlerOption {
	return func(meta *RouteMeta) {
		meta.QueryParams = append(meta.QueryParams, params...)
	}
}

// --- Typed path params (auto parse) ---

// PathParam is a typed wrapper around a path parameter.
//
// Example:
//
//	func(w http.ResponseWriter, r *http.Request, id openapi.PathParam[int]) (...)
//
// Use id.Value when present.
type PathParam[T any] struct {
	Name  string
	Value T
	OK    bool
}

// PathParamValue reads and parses a path parameter into a concrete type.
func PathParamValue[T any](r *http.Request, name string) (PathParam[T], error) {
	raw := PathValue(r, name)
	if raw == "" {
		var z T
		return PathParam[T]{Name: name, Value: z, OK: false}, errors.New("missing path parameter: " + name)
	}

	v, err := parsePrimitive[T](raw)
	if err != nil {
		var z T
		return PathParam[T]{Name: name, Value: z, OK: false}, err
	}

	return PathParam[T]{Name: name, Value: v, OK: true}, nil
}

// parsePrimitive converts a string into a typed primitive (string/int/float/bool).
func parsePrimitive[T any](raw string) (T, error) {
	var out T
	t := reflect.TypeOf(out)
	if t == nil {
		return out, errors.New("invalid type")
	}

	// Support aliases by checking kind.
	switch t.Kind() {
	case reflect.String:
		return any(raw).(T), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return out, err
		}
		v := reflect.New(t).Elem()
		v.SetInt(i)
		return v.Interface().(T), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return out, err
		}
		v := reflect.New(t).Elem()
		v.SetUint(u)
		return v.Interface().(T), nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return out, err
		}
		v := reflect.New(t).Elem()
		v.SetFloat(f)
		return v.Interface().(T), nil
	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return out, err
		}
		v := reflect.New(t).Elem()
		v.SetBool(b)
		return v.Interface().(T), nil
	default:
		return out, errors.New("unsupported param type: " + t.String())
	}
}

func openapiTypeToSchemaType(t ParamType) *openapi3.Types {
	switch t {
	case ParamInteger:
		return &openapi3.Types{"integer"}
	case ParamNumber:
		return &openapi3.Types{"number"}
	case ParamBoolean:
		return &openapi3.Types{"boolean"}
	default:
		return &openapi3.Types{"string"}
	}
}

// QueryValue reads a query parameter from URL and parses it into the requested type.
func QueryValue[T any](r *http.Request, name string) (T, bool, error) {
	var z T
	if r == nil || r.URL == nil {
		return z, false, errors.New("nil request")
	}
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return z, false, nil
	}
	v, err := parsePrimitive[T](raw)
	if err != nil {
		return z, false, err
	}
	return v, true, nil
}

// QueryValues reads repeated query params (?id=1&id=2) and parses into []T.
func QueryValues[T any](r *http.Request, name string) ([]T, bool, error) {
	if r == nil || r.URL == nil {
		return nil, false, errors.New("nil request")
	}
	vals, ok := r.URL.Query()[name]
	if !ok || len(vals) == 0 {
		return nil, false, nil
	}
	out := make([]T, 0, len(vals))
	for _, raw := range vals {
		v, err := parsePrimitive[T](raw)
		if err != nil {
			return nil, false, err
		}
		out = append(out, v)
	}
	return out, true, nil
}

// Helper to add query params into operation.
func addQueryParams(op *openapi3.Operation, qps []QueryParam) {
	for _, qp := range qps {
		if strings.TrimSpace(qp.Name) == "" {
			continue
		}
		p := &openapi3.Parameter{
			Name:        qp.Name,
			In:          openapi3.ParameterInQuery,
			Required:    qp.Required,
			Description: qp.Description,
			Schema:      &openapi3.SchemaRef{Value: &openapi3.Schema{Type: openapiTypeToSchemaType(qp.Type)}},
		}
		op.AddParameter(p)
	}
}

// Utility for tests/examples.
func withQuery(r *http.Request, values url.Values) *http.Request {
	if r.URL == nil {
		return r
	}
	r2 := *r
	u := *r.URL
	u.RawQuery = values.Encode()
	r2.URL = &u
	return &r2
}

type PathParamSpec struct {
	Name        string
	Type        ParamType
	Required    bool
	Description string
}

func WithPathParam(name string, typ ParamType, required bool, description string) HandlerOption {
	return func(meta *RouteMeta) {
		meta.PathParams = append(meta.PathParams, PathParamSpec{
			Name:        name,
			Type:        typ,
			Required:    required,
			Description: description,
		})
	}
}
