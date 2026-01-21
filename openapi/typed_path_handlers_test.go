package openapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestGETP1T_ParsesAndDocumentsParamType(t *testing.T) {
	r := NewRouter()

	GETP1T[struct{}, map[string]any, int](r, "/users/{id}", "id", func(w http.ResponseWriter, req *http.Request, _ struct{}, id Param[int]) (map[string]any, int, error) {
		_ = w
		_ = req
		return map[string]any{"id": id.Value}, http.StatusOK, nil
	})

	doc := BuildSpec(r.Routes(), Config{Title: "T", Version: "1"})
	p := doc.Paths.Find("/users/{id}")
	if p == nil || p.Get == nil {
		t.Fatalf("expected GET /users/{id}")
	}

	// Ensure documented as integer
	found := false
	for _, pr := range p.Get.Parameters {
		if pr.Value == nil {
			continue
		}
		if pr.Value.In != openapi3.ParameterInPath || pr.Value.Name != "id" {
			continue
		}
		found = true
		if pr.Value.Schema == nil || pr.Value.Schema.Value == nil || pr.Value.Schema.Value.Type == nil {
			t.Fatalf("expected schema type")
		}
		if len(*pr.Value.Schema.Value.Type) == 0 || (*pr.Value.Schema.Value.Type)[0] != "integer" {
			t.Fatalf("expected integer type, got %#v", pr.Value.Schema.Value.Type)
		}
	}
	if !found {
		t.Fatalf("expected path param id")
	}

	// Ensure runtime parsing works
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
