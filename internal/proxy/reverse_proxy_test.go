package proxy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"apigateway/internal/model"
)

func testFactory(route *model.Route) LoadBalancer {
	return NewRoundRobinBalancer()
}

func TestReverseProxyHandler_SuccessfulProxy(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer upstream.Close()

	route := &model.Route{
		Path:      "/api/test",
		Method:    "GET",
		Upstreams: []model.Upstream{{Address: upstream.URL, Weight: 1}},
	}

	matcher := func(path, method string) (*model.Route, map[string]string, bool) {
		return route, nil, true
	}

	handler := NewReverseProxyHandler(matcher, testFactory)
	req := httptest.NewRequest("GET", "/api/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}

func TestReverseProxyHandler_NoRoute(t *testing.T) {
	matcher := func(path, method string) (*model.Route, map[string]string, bool) {
		return nil, nil, false
	}

	handler := NewReverseProxyHandler(matcher, testFactory)
	req := httptest.NewRequest("GET", "/unknown", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != 404 {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestReverseProxyHandler_UpstreamDown(t *testing.T) {
	route := &model.Route{
		Path:      "/api/test",
		Method:    "GET",
		Upstreams: []model.Upstream{{Address: "http://localhost:1", Weight: 1}},
	}

	matcher := func(path, method string) (*model.Route, map[string]string, bool) {
		return route, nil, true
	}

	handler := NewReverseProxyHandler(matcher, testFactory)
	req := httptest.NewRequest("GET", "/api/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != 502 {
		t.Errorf("status = %d, want 502", rec.Code)
	}
}

func TestReverseProxyHandler_ForwardedHeaders(t *testing.T) {
	var gotXFF, gotXRealIP, gotXFP string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotXFF = r.Header.Get("X-Forwarded-For")
		gotXRealIP = r.Header.Get("X-Real-IP")
		gotXFP = r.Header.Get("X-Forwarded-Proto")
		w.WriteHeader(200)
	}))
	defer upstream.Close()

	route := &model.Route{
		Path:      "/api/test",
		Method:    "GET",
		Upstreams: []model.Upstream{{Address: upstream.URL, Weight: 1}},
	}

	matcher := func(path, method string) (*model.Route, map[string]string, bool) {
		return route, nil, true
	}

	handler := NewReverseProxyHandler(matcher, testFactory)
	req := httptest.NewRequest("GET", "/api/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if gotXFF == "" {
		t.Error("X-Forwarded-For should be set")
	}
	if gotXRealIP == "" {
		t.Error("X-Real-IP should be set")
	}
	if gotXFP != "http" {
		t.Errorf("X-Forwarded-Proto = %q, want http", gotXFP)
	}
}
