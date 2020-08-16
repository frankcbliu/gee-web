package gee

import (
	"reflect"
	"testing"
)

func Test_parsePattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    []string
	}{ // 单测用例
		{name: "/p/*", pattern: "/p/*", want: []string{"p", "*"}},
		{name: "/p/:name", pattern: "/p/:name", want: []string{"p", ":name"}},
		{name: "/p/*name/*", pattern: "/p/*name/*", want: []string{"p", "*name"}},
		{name: "/p/:name/b/*", pattern: "/p/:name/b/*", want: []string{"p", ":name", "b", "*"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePattern(tt.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePattern() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func newTestRouter() *router {
	r := NewRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}
func Test_router_getRoute(t *testing.T) {
	type args struct {
		method string
		path   string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 map[string]string
	}{ // 单测用例
		{name: "geek", args: args{"GET", "/hello/geek"}, want: "/hello/:name", want1: map[string]string{"name": "geek"}},
		{name: "frank", args: args{"GET", "/hello/frank"}, want: "/hello/:name", want1: map[string]string{"name": "frank"}},
		{name: "hello_b_c", args: args{"GET", "/hello/b/c"}, want: "/hello/b/c", want1: map[string]string{}},
		{name: "assets", args: args{"GET", "/assets/233.jpg"}, want: "/assets/*filepath", want1: map[string]string{"filepath": "233.jpg"}},
	}
	r := newTestRouter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := r.getRoute(tt.args.method, tt.args.path)
			if !reflect.DeepEqual(got.pattern, tt.want) {
				t.Errorf("getRoute() got = %v, want %v", got.pattern, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getRoute() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
