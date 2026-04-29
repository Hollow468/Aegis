package router

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"apigateway/internal/model"
)

// RouteMatch holds the matched route and extracted path parameters.
type RouteMatch struct {
	Route  *model.Route
	Params map[string]string
}

// segmentType classifies path segments.
type segmentType int

const (
	segExact    segmentType = iota // literal segment
	segParam                       // :name
	segRegex                       // {name:pattern}
	segWildcard                    // *
)

// trieNode represents a node in the routing trie.
type trieNode struct {
	children     map[string]*trieNode // exact segment children
	paramChild   *trieNode            // :param child
	paramName    string               // param name (without colon)
	regexChild   *trieNode            // {param:regex} child
	regex        *regexp.Regexp       // compiled regex
	regexName    string               // param name for regex
	wildcardChild *trieNode           // * child
	methodRoutes map[string]*model.Route
	isLeaf       bool
}

func newTrieNode() *trieNode {
	return &trieNode{
		children:     make(map[string]*trieNode),
		methodRoutes: make(map[string]*model.Route),
	}
}

// TrieRouter implements a trie-based routing engine supporting exact, prefix, and regex matching.
type TrieRouter struct {
	root *trieNode
	mu   sync.RWMutex
}

// NewTrieRouter creates a new TrieRouter.
func NewTrieRouter() *TrieRouter {
	return &TrieRouter{root: newTrieNode()}
}

// Insert adds a route to the trie for the given path and HTTP method.
func (tr *TrieRouter) Insert(path, method string, route *model.Route) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	segments := parsePath(path)
	node := tr.root

	for _, seg := range segments {
		segType, name, pattern := classifySegment(seg)

		switch segType {
		case segExact:
			child, ok := node.children[seg]
			if !ok {
				child = newTrieNode()
				node.children[seg] = child
			}
			node = child

		case segParam:
			if node.paramChild == nil {
				node.paramChild = newTrieNode()
				node.paramName = name
			}
			node = node.paramChild

		case segRegex:
			re, err := regexp.Compile("^" + pattern + "$")
			if err != nil {
				return fmt.Errorf("invalid regex in segment %s: %w", seg, err)
			}
			if node.regexChild == nil {
				node.regexChild = newTrieNode()
				node.regex = re
				node.regexName = name
			} else {
				// Update regex for existing child (multiple methods on same path)
				node.regex = re
				node.regexName = name
			}
			node = node.regexChild

		case segWildcard:
			if node.wildcardChild == nil {
				node.wildcardChild = newTrieNode()
			}
			node = node.wildcardChild
		}
	}

	node.isLeaf = true
	node.methodRoutes[method] = route
	return nil
}

// Match finds a route matching the given path and HTTP method.
// Priority: exact > regex > param > wildcard.
func (tr *TrieRouter) Match(path, method string) (*RouteMatch, bool) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	segments := parsePath(path)
	params := make(map[string]string)

	node := tr.matchNode(tr.root, segments, 0, params)
	if node == nil {
		return nil, false
	}

	route, ok := node.methodRoutes[method]
	if !ok {
		return nil, false
	}

	// Copy params
	p := make(map[string]string, len(params))
	for k, v := range params {
		p[k] = v
	}

	return &RouteMatch{Route: route, Params: p}, true
}

func (tr *TrieRouter) matchNode(node *trieNode, segments []string, idx int, params map[string]string) *trieNode {
	if idx == len(segments) {
		if node.isLeaf {
			return node
		}
		return nil
	}

	seg := segments[idx]

	// 1. Exact match (highest priority)
	if child, ok := node.children[seg]; ok {
		if result := tr.matchNode(child, segments, idx+1, params); result != nil {
			return result
		}
	}

	// 2. Regex match
	if node.regexChild != nil && node.regex != nil && node.regex.MatchString(seg) {
		params[node.regexName] = seg
		if result := tr.matchNode(node.regexChild, segments, idx+1, params); result != nil {
			return result
		}
		delete(params, node.regexName)
	}

	// 3. Param match
	if node.paramChild != nil {
		params[node.paramName] = seg
		if result := tr.matchNode(node.paramChild, segments, idx+1, params); result != nil {
			return result
		}
		delete(params, node.paramName)
	}

	// 4. Wildcard match (lowest priority, matches everything remaining)
	if node.wildcardChild != nil {
		return node.wildcardChild
	}

	return nil
}

// classifySegment determines the type and extracts name/pattern.
func classifySegment(seg string) (segType segmentType, name, pattern string) {
	if seg == "*" {
		return segWildcard, "", ""
	}

	// Curly brace segment: {name} or {name:pattern}
	if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
		inner := seg[1 : len(seg)-1]
		parts := strings.SplitN(inner, ":", 2)
		if len(parts) == 2 && parts[1] != "" {
			// Has regex pattern: {name:pattern}
			return segRegex, parts[0], parts[1]
		}
		// No pattern: {name} — treat as param
		return segParam, inner, ""
	}

	// Param: :name
	if strings.HasPrefix(seg, ":") {
		return segParam, seg[1:], ""
	}

	// Exact
	return segExact, "", ""
}

// parsePath splits a path into segments, removing empty segments.
func parsePath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}

// Routes returns all registered routes (for debugging).
func (tr *TrieRouter) Routes() []string {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	var routes []string
	tr.collectRoutes(tr.root, "", &routes)
	return routes
}

func (tr *TrieRouter) collectRoutes(node *trieNode, prefix string, routes *[]string) {
	if node.isLeaf {
		for method, route := range node.methodRoutes {
			*routes = append(*routes, fmt.Sprintf("%s %s -> %v", method, prefix, route.Upstreams))
		}
	}
	for seg, child := range node.children {
		tr.collectRoutes(child, prefix+"/"+seg, routes)
	}
	if node.paramChild != nil {
		tr.collectRoutes(node.paramChild, prefix+"/:"+node.paramName, routes)
	}
	if node.regexChild != nil {
		tr.collectRoutes(node.regexChild, prefix+"/{"+node.regexName+":...}", routes)
	}
	if node.wildcardChild != nil {
		tr.collectRoutes(node.wildcardChild, prefix+"/*", routes)
	}
}
