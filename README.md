# TengoHTTP
A simple [Tengo](https://github.com/d5/tengo) HTTP-server library, used in [ServeMSX](https://github.com/damiva/ServeMSX) project.
## Usage in GOLANG
### Installation
```
go get github.com/damiva/TengoHTTP
```
### Function
```golang
// GetModuleMAP returns Builtin Tengo Module, where:
// c - HTTP Client for HTTP Request, if nil - &http.Client{} is used
// vars - custom Tengo variables added to the module
func GetModuleMAP(w http.ResponseWriter, r *http.Request, c *http.Client, vars map[string]tengo.Object) map[string]tengo.Object
```
### Example
```golang
package main

import (
	"log"
	"net/http"

	tengo "github.com/d5/tengo/v2"
	tengohttp "github.com/damiva/TengoHTTP"
)

const Name, Version = "MyServer", "1.00"

var script = []byte(`
srv := import("server")
rsp := srv.request("http://google.com", "HEAD")
srv.log("Request answer:", is_error(rsp) ? rsp : rsp.status)
srv.write("Hello from ", srv.name, " v. ", srv.version)
`)

func tengoLog(args ...tengo.Object) (tengo.Object, error) {
	if len(arg) == 0 {
		return nil, tengo.ErrWrongNumArguments
	}
	var v []interface{}
	for _, a := range args {
		if s, _ := tengo.ToString(a); s != "" {
			v = append(v, s)
		}
	}
	log.Println(v...)
	return nil, nil
}
func main() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		tng := tengo.NewScript(script)
		mod := new(tengo.ModuleMap)
		mod.AddBuiltinModule("server", tengohttp.GetModuleMAP(w, r, nil, map[string]tengo.Object{
			"name":    &tengo.String{Value: Name},
			"version": &tengo.String{Value: Version},
			"log":     &tengo.UserFunction{Name: "log", Value: tengoLog},
		}))
		tng.SetImports(mod)
		log.Println("Run test:")
		if _, e := tng.Run(); e != nil {
			log.Println("Error:", e)
		} else {
			log.Println("Done")
		}
	})
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
```
## Usage in Tengo
```golang
srv := import("server")
```
### Variables
- `proto {string}`: the protocol version for the request;
- `method {string}`: the HTTP method (GET, POST, PUT, etc.) of the request;
- `host {string}`: the host on which the URL is sought;
- `remote_addr {string}`: the network address that sent the request;
- `header {map of arrays of strings}`: the request header fields of the request;
- `uri {string}`: unmodified request-target of the Request-Line (RFC 7230, Section 3.1.1) as sent by the client to a server;
- `raw_query {string}`: encoded URL query string (without '?') of the request.

### Functions
- `read([string/bool]) => {bytes/string/map/error}`: where:
	- if the argument is absent: returns the body of the request as {bytes}
	- if the argument is {string}: returns the first value for the named component of the query as {string}, POST and PUT body parameters take precedence over URL query string values;
	- if the argument = *true*: returns fields of the POST/PUT query as {map of arrays of strings};
	- if the argument = *false*: returns fileds of the GET query (URL query string values) as {map of arrays of strings};
- `write([any]...) => [error]`:
- `request({string}[, string/bytes/map]) => {map/error}`:
- `uri_encode({string}[, bool]) => {string}`:
- `uri_decode({string}[, bool]) => {string/error}`:
- `url_parse([string/map]) => {map/string/error}`:
- `url_resolve({string}[, string]) => {string/error}`:
- `query_parse({string/map}) => {map/string/error}`:

