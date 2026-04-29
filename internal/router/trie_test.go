package router

import (
	"testing"

	"apigateway/internal/model"
)

func TestTrieRouter_ExactMatch(t *testing.T) {
	tr := NewTrieRouter()
	route := &model.Route{Path: "/api/users", Method: "GET"}
	tr.Insert("/api/users", "GET", route)

	tests := []struct {
		path    string
		method  string
		wantOK  bool
	}{
		{"/api/users", "GET", true},
		{"/api/users", "POST", false},
		{"/api/other", "GET", false},
		{"/api/users/123", "GET", false},
	}

	for _, tt := range tests {
		rm, ok := tr.Match(tt.path, tt.method)
		if ok != tt.wantOK {
			t.Errorf("Match(%q, %q) = %v, want %v", tt.path, tt.method, ok, tt.wantOK)
		}
		if ok && rm.Route != route {
			t.Errorf("Match(%q, %q) returned wrong route", tt.path, tt.method)
		}
	}
}

func TestTrieRouter_PrefixMatch(t *testing.T) {
	tr := NewTrieRouter()
	route := &model.Route{Path: "/api/*", Method: "GET"}
	tr.Insert("/api/*", "GET", route)

	tests := []struct {
		path   string
		method string
		wantOK bool
	}{
		{"/api/files", "GET", true},
		{"/api/anything/here", "GET", true},
		{"/other/path", "GET", false},
	}

	for _, tt := range tests {
		_, ok := tr.Match(tt.path, tt.method)
		if ok != tt.wantOK {
			t.Errorf("Match(%q, %q) = %v, want %v", tt.path, tt.method, ok, tt.wantOK)
		}
	}
}

func TestTrieRouter_RegexMatch(t *testing.T) {
	tr := NewTrieRouter()
	route := &model.Route{Path: "/api/users/{id:[0-9]+}", Method: "GET"}
	tr.Insert("/api/users/{id:[0-9]+}", "GET", route)

	tests := []struct {
		path    string
		method  string
		wantOK  bool
		wantID  string
	}{
		{"/api/users/123", "GET", true, "123"},
		{"/api/users/abc", "GET", false, ""},
		{"/api/users/0", "GET", true, "0"},
		{"/api/users", "GET", false, ""},
	}

	for _, tt := range tests {
		rm, ok := tr.Match(tt.path, tt.method)
		if ok != tt.wantOK {
			t.Errorf("Match(%q, %q) = %v, want %v", tt.path, tt.method, ok, tt.wantOK)
			continue
		}
		if ok && tt.wantID != "" {
			if rm.Params["id"] != tt.wantID {
				t.Errorf("Match(%q) param id = %q, want %q", tt.path, rm.Params["id"], tt.wantID)
			}
		}
	}
}

func TestTrieRouter_Priority(t *testing.T) {
	tr := NewTrieRouter()
	exact := &model.Route{Path: "/api/users", Method: "GET"}
	param := &model.Route{Path: "/api/users/{id}", Method: "GET"}
	wildcard := &model.Route{Path: "/api/*", Method: "GET"}

	tr.Insert("/api/users", "GET", exact)
	tr.Insert("/api/users/{id}", "GET", param)
	tr.Insert("/api/*", "GET", wildcard)

	// Exact should win
	rm, ok := tr.Match("/api/users", "GET")
	if !ok || rm.Route != exact {
		t.Error("exact match should have highest priority")
	}

	// Param should win over wildcard
	rm, ok = tr.Match("/api/users/123", "GET")
	if !ok || rm.Route != param {
		t.Error("param match should win over wildcard")
	}

	// Wildcard as fallback
	rm, ok = tr.Match("/api/other", "GET")
	if !ok || rm.Route != wildcard {
		t.Error("wildcard should match as fallback")
	}
}

func TestTrieRouter_MethodSeparation(t *testing.T) {
	tr := NewTrieRouter()
	getRoute := &model.Route{Path: "/api/users", Method: "GET"}
	postRoute := &model.Route{Path: "/api/users", Method: "POST"}

	tr.Insert("/api/users", "GET", getRoute)
	tr.Insert("/api/users", "POST", postRoute)

	rm, ok := tr.Match("/api/users", "GET")
	if !ok || rm.Route != getRoute {
		t.Error("GET route should match GET method")
	}

	rm, ok = tr.Match("/api/users", "POST")
	if !ok || rm.Route != postRoute {
		t.Error("POST route should match POST method")
	}
}

func TestTrieRouter_RootPath(t *testing.T) {
	tr := NewTrieRouter()
	route := &model.Route{Path: "/", Method: "GET"}
	tr.Insert("/", "GET", route)

	rm, ok := tr.Match("/", "GET")
	if !ok || rm.Route != route {
		t.Error("root path should match")
	}
}

func TestTrieRouter_EmptyPath(t *testing.T) {
	tr := NewTrieRouter()
	_, ok := tr.Match("", "GET")
	if ok {
		t.Error("empty path should not match anything")
	}
}

func TestTrieRouter_MultipleParams(t *testing.T) {
	tr := NewTrieRouter()
	route := &model.Route{Path: "/api/users/{userId}/posts/{postId}", Method: "GET"}
	tr.Insert("/api/users/{userId}/posts/{postId}", "GET", route)

	rm, ok := tr.Match("/api/users/42/posts/99", "GET")
	if !ok {
		t.Fatal("should match multi-param route")
	}
	if rm.Params["userId"] != "42" {
		t.Errorf("userId = %q, want %q", rm.Params["userId"], "42")
	}
	if rm.Params["postId"] != "99" {
		t.Errorf("postId = %q, want %q", rm.Params["postId"], "99")
	}
}
