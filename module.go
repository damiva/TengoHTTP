package tengohttp

import (
	"net/http"

	"github.com/d5/tengo/v2"
)

type server struct {
	w http.ResponseWriter
	r *http.Request
	c *http.Client
	h bool
}

// GetModuleMAP returns Builtin Tengo Module, where:
// cln  - HTTP Client for HTTP Request, if nil - &http.Client{} is used
// vars - custom variables can be added to the module, custom variables take precedence over builtin variables
func GetModuleMAP(w http.ResponseWriter, r *http.Request, cln *http.Client, vars map[string]tengo.Object) map[string]tengo.Object {
	if cln == nil {
		cln = &http.Client{}
	}
	s := &server{w, r, cln, false}
	ret := map[string]tengo.Object{
		"proto":       &tengo.String{Value: r.Proto},
		"method":      &tengo.String{Value: r.Method},
		"host":        &tengo.String{Value: r.Host},
		"remote_addr": &tengo.String{Value: r.RemoteAddr},
		"header":      vals2map(r.Header),
		"uri":         &tengo.String{Value: r.RequestURI},
		"write":       &tengo.UserFunction{Name: "write", Value: s.write},
		"read":        &tengo.UserFunction{Name: "read", Value: s.read},
		"request":     &tengo.UserFunction{Name: "request", Value: s.request},
		"encode_uri":  &tengo.UserFunction{Name: "encode_uri", Value: encode},
		"decode_uri":  &tengo.UserFunction{Name: "decode_uri", Value: decode},
		"parse_url":   &tengo.UserFunction{Name: "parse_url", Value: s.parse},
		"resolve_url": &tengo.UserFunction{Name: "resolve_url", Value: s.resolve},
	}
	for k, v := range vars {
		ret[k] = v
	}
	return ret
}
