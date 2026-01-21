package openapi

import (
	"net/http"
	"reflect"
)

// TypedHandler is a convenience handler signature that enables full auto schema
// inference via type parameters.
//
// TReq is inferred as request body schema (unless it is struct{}).
// TRes is inferred as response schema (unless it is struct{}).
//
// Your handler returns (value, statusCode, error). If error != nil, the default
// behavior is to respond with 500 and the error string.
//
// If you need more control, you can still use the classic net/http handler and
// optionally attach schemas via WithRequestSchema/WithResponseSchema.
type TypedHandler[TReq any, TRes any] func(w http.ResponseWriter, r *http.Request, req TReq) (res TRes, status int, err error)

// isZeroStructType returns true when T is exactly struct{}.
func isZeroStructType[T any]() bool {
	var zero T
	t := reflect.TypeOf(zero)
	return t != nil && t.Kind() == reflect.Struct && t.NumField() == 0
}

// wrapTyped converts a typed handler into http.HandlerFunc.
func wrapTyped[TReq any, TRes any](h TypedHandler[TReq, TRes]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqVal TReq
		if !isZeroStructType[TReq]() {
			// ignore decode error here; handler can validate/return custom errors later.
			_ = Bind(r, &reqVal)
		}

		res, code, err := h(w, r, reqVal)
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
	}
}

func typedOptions[TReq any, TRes any]() (reqOpt, resOpt HandlerOption) {
	var reqZero TReq
	var resZero TRes

	if !isZeroStructType[TReq]() {
		reqOpt = WithRequestSchema(reqZero)
	}
	if !isZeroStructType[TRes]() {
		resOpt = WithResponseSchema(resZero)
	}

	return reqOpt, resOpt
}

func mergeOpts(base []HandlerOption, add ...HandlerOption) []HandlerOption {
	out := make([]HandlerOption, 0, len(base)+len(add))
	out = append(out, base...)
	out = append(out, add...)
	return out
}

// GETT registers a typed GET handler. Request schema is inferred from TReq and response from TRes.
func GETT[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	reqOpt, resOpt := typedOptions[TReq, TRes]()
	base := make([]HandlerOption, 0, 2)
	if reqOpt != nil {
		base = append(base, reqOpt)
	}
	if resOpt != nil {
		base = append(base, resOpt)
	}
	r.GET(path, wrapTyped(h), mergeOpts(base, opts...)...)
}

// POSTT registers a typed POST handler.
func POSTT[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	reqOpt, resOpt := typedOptions[TReq, TRes]()
	base := make([]HandlerOption, 0, 2)
	if reqOpt != nil {
		base = append(base, reqOpt)
	}
	if resOpt != nil {
		base = append(base, resOpt)
	}
	r.POST(path, wrapTyped(h), mergeOpts(base, opts...)...)
}

// PUTT registers a typed PUT handler.
func PUTT[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	reqOpt, resOpt := typedOptions[TReq, TRes]()
	base := make([]HandlerOption, 0, 2)
	if reqOpt != nil {
		base = append(base, reqOpt)
	}
	if resOpt != nil {
		base = append(base, resOpt)
	}
	r.PUT(path, wrapTyped(h), mergeOpts(base, opts...)...)
}

// DELETET registers a typed DELETE handler.
func DELETET[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	reqOpt, resOpt := typedOptions[TReq, TRes]()
	base := make([]HandlerOption, 0, 2)
	if reqOpt != nil {
		base = append(base, reqOpt)
	}
	if resOpt != nil {
		base = append(base, resOpt)
	}
	r.DELETE(path, wrapTyped(h), mergeOpts(base, opts...)...)
}

// PATCHT registers a typed PATCH handler.
func PATCHT[TReq any, TRes any](r *Router, path string, h TypedHandler[TReq, TRes], opts ...HandlerOption) {
	reqOpt, resOpt := typedOptions[TReq, TRes]()
	base := make([]HandlerOption, 0, 2)
	if reqOpt != nil {
		base = append(base, reqOpt)
	}
	if resOpt != nil {
		base = append(base, resOpt)
	}
	r.PATCH(path, wrapTyped(h), mergeOpts(base, opts...)...)
}
