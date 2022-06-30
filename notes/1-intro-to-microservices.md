# Tooling and Software
- Nic is using a windows machine!
- Windows has a Unix based terminal now
- apparently a little more compatible
- that being said, it doesn't matter which OS you're using when writing microservices in Go
- Extensions
  - Go language support
  - Docker
  - Live Share (won't need it for purposes of this)
- He's running `go1.13.5`
  - I'm on `go1.18.3`
- Nic loves Docker

# Your first service using the Go standard library
- What does it take to build a service?
- Nothing new yet. Notes will likely be sparse.
- Oh sike immediately making a web server
  - [http.ListenAndServe](https://pkg.go.dev/net/http@go1.18.3#ListenAndServe)
  - `func ListenAndServe(addr string, handler Handler) error`
```go
package main

import "net/http"

func main() {
	// bind to every IP at port 9090
	http.ListenAndServe(":9090", nil)
}
```
- Ignoring the handler for now
- this server is listening at _all_ IP addresses and port `9090`
- This is a webserver!
- We can see that it does in fact work but nothing is being sent since we haven't set up the logic to handle requests
```bash
$ curl -v localhost:9090
*   Trying 127.0.0.1:9090...
* Connected to localhost (127.0.0.1) port 9090 (#0)
> GET / HTTP/1.1
> Host: localhost:9090
> User-Agent: curl/7.79.1
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 404 Not Found
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Thu, 30 Jun 2022 18:45:11 GMT
< Content-Length: 19
< 
404 page not found
* Connection #0 to host localhost left intact
```
- To handle a request, we're using [http.HandleFunc](https://pkg.go.dev/net/http@go1.18.3#HandleFunc)
  - `func HandleFunc(pattern string, handler func(ResponseWriter, *Request))`
```go
func main() {
	http.HandleFunc("/", func(http.ResponseWriter, *http.Request) {
		log.Println(("Hello World"))
	})

	http.ListenAndServe(":9090", nil)
}
```
- now we'll see `Hello World` logged _in the server_
  - mike: look into `log.Println`
- `HandleFunc` registers a function to a path on a thing called the "default serve mux" (omg yes a multiplexer)
- the DSM is an HTTP handler
- When we don't specify a handler in `ListenAndServe`, Go defaults to the DSM
- serve mux is responsible for redirecting paths
- `ServeMux` is an HTTP handler
- you can create your own serve mux
- a serve mux has a register method, which registers a handler at a path
- a handler is an interface that implements `ServeHTTP(ResponseWriter, *Request)`
- So, looking back, the `HandleFunc` method registers this function literal with the path provided (`"/"`)
- Wow it's literally a mux
- okay let's register another route
```go
http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
	log.Println(("Goodbye World"))
})
```
- the ServeMux does greedy matching, so anything _other than_ `/goodbye` will match the root path
- How do we read and write to the request?
- Using the response writer and the http request
- "A ResponseWriter interface is used by an HTTP handler to construct an HTTP response."
- The request is the... request. Has metadata and path params etc.
- The request has a method called `Body`, which implements the `io.ReadCloser` interface (I know that reference!)
- We can now implement whatever method to read that data, like `ioutil.ReadAll`
  - [ioutil.ReadAll](https://pkg.go.dev/io/ioutil@go1.18.3#ReadAll)
  - `func ReadAll(r io.Reader) ([]byte, error)`
  - Oh interesting... It returns a byte slice but we can just log it normally
```go
http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
	log.Println(("Hello World"))
	d, _ := ioutil.ReadAll(r.Body)

	log.Printf("Data %s\n", d)
})
```
- This just logs the data in the server though
- we want to write the data back to the user!
- Response writer has a `Write` method and implements the `io.Writer` interface! So we can use it anywhere that accepts a Writer, like `Fprintf`
```go
http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
	log.Println(("Hello World"))
	d, _ := ioutil.ReadAll(r.Body)

	fmt.Fprintf(rw, "Hello %s", d)
})
```
- Now it get's returned to the user!
- Okay Mike's turn to reflect
  - `http.ResponseWriter` is an interface with a write method. That means that whatever the Go http server passes into `rw` will have some write method that implements the expected signature
  - other methods that operate on the writer don't care what that "write" looks like, so long as the contract is maintained. 
  - `Fprintf` says: _Gimme a string and a way to writer to feed it into_
  - WE say: _here's a response writer and a string_
  - and finally the response writer doesn't say anything but it thanks `Fprintf` for the string and does the "write", which is sending the data to the user
  - Okay this is cool. It kind of lets us decide _how_ the response writer gets used
- Nic adds this is actually a very performant server
- WE can however use the build in methods in `rw` like `WriteHeader`, which lets us set the response status
- OMG Go has all the status codes already build in as constants
- We _could_ write an error message like
```go
if err != nil {
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte("Oops"))
	return
}
```
- ... but go has a convenience method for us
  - [http.Error](https://pkg.go.dev/net/http@go1.18.3#Error)
  - `func Error(w ResponseWriter, error string, code int)`
- I mean... wow. Just wow.
```go
if err != nil {
	http.Error(rw, "Oops", http.StatusBadRequest)
	return
}
```
- we still need the `return` statement since `Error` doesn't terminate the code
- All together now
```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println(("Hello World"))
		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, "Oops", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(rw, "Hello %s", d)
	})

	http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
		log.Println(("Goodbye World"))
	})

	http.ListenAndServe(":9090", nil)
}

```
- If we want to think about testability and DX, we need to think more about code structure
