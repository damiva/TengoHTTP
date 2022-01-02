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
// cln  - HTTP Client for HTTP Request, if nil - &http.Client{} is used
// vars - custom variables can be added to the module, custom variables take precedence over builtin variables
func GetModuleMAP(w http.ResponseWriter, r *http.Request, cln *http.Client, vars map[string]tengo.Object) map[string]tengo.Object
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
- `header {map of arrays of strings}`: the request header fields;
- `uri {string}`: unmodified request-target of the Request-Line (RFC 7230, Section 3.1.1) as sent by the client to a server;

### Functions
- `read() => {bytes/error}`: returns the request body as is;
- `read({bool}) => {map of arrays of strings}`: parses the request and:
	- if {bool} is falsy, returns the fields of the query, POST/PUT body parameters take precedence over URL query string values;
	- else returns the fields of the POST/PUT body parameters;
- `read({bool}, {string}) => {string/undefined}`: parses the request and returns the first value of the named {string} field of:
	- if {bool} is falsy: the query, POST/PUT body parameters take precedence over URL query string values;(
	- else: the POST/PUT body parameters;
- `write([any]...) => {undefined/error}`: writes the response, where:
	- if an argument is {map} sets the headers of the response, it should be map of arrays of strings, and should be set before status or body writings;
	- if an argument is {int} set the http status of the response to the {int} code, it should be set before body writings, if {int} between 300 & 399 and the next argument is {string} sends the response as redirect with the code = {int} and the url = {string};
	- any other arguments are writen to the body of the response;
- `request({string}[, string/bytes/map]) => {map/error}`: does the http request and returns the answer, where:
	- first argument is the url of the request;
	- if the second argument is {string} sets the http method (default is *GET*);
	- if the second argument is {bytes} sends the request body with http method *POST*;
	- if the second argument is map, sends the request with the following parameters in the map:
		- "body" {map/bytes/string}: the body of the request (if {map} it will be encoded as form);
		- "method" {string}: the http method (if "body" absent default is *GET*, else *POST*);
		- "query" {map/string}: the uri query (if {map} it will be encoded as form);
		- "header" {map}: the header fileds of the request (map of strings/array of strings);
		- "cookies" {array};
		- "user" {string};
		- "pass" {string};
		- "follow" {bool}: if it is falsy it does not follow the redirects (default is follows up to 10 times);
		- "timeout" {int}: in seconds;
	- the answer contains the response as a map of:
	 	- "status" {int}: status code;
	 	- "user" {string};
	 	- "pass" {string};
	 	- "header" {map};
	 	- "cookies" {array};
	 	- "body" {bytes};
	 	- "size" {int}: number of bytes;
	 	- "url" {string}: the final url (after all redirects);
- `encode({string/map}[, bool]) => {string}`:
- `decode({string}[, bool]) => {string/error}`:
- `parse([string]) => {map/string/error}`:
- `resolve({string}[, string]) => {string/error}`:

