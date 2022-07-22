# Intro
* Nic feels that when using the std lib, you end up writing a lot of boilerplate code
* Now we're going to refactor again using the [Gorilla web toolkit](https://www.gorillatoolkit.org/).

# Gorilla
* includes a few different packages, including the [mux package](https://github.com/gorilla/mux)
* This improves on the std mux
  * placeholders
  * define path IDs with a regex
* handles middleware!

# Getting to it now
* First thing's first, replace the servemux
* Now, Mike, let's remind ourselves _what_ the servemux _is_
  * the sm is the request router! It's where you define what to do at a given request path
  * `/`: go here
  * `/dogs`: go there
  * ok
* Need to refactor the routers too
* To install Gorilla, run

```bash
$ go get github.com/gorilla/mux
```

* This is the analog to `npm install <pkg_name>` in Node
* the gorilla mux has sub-router functionality
* So there's a few things going on here
  1. you can define a route specifically for a certain method, as in "only trigger this on `GET /api/route`
  2. you can make a subrouter, which means something that I don't understand yet. From the docs
> Routes can be used as subrouters: nested routes are only tested if the parent route matches. This is useful to define groups of routes that share common conditions like a host, a path prefix or other repeated attributes. As a bonus, this optimizes request matching.

* Making something a subrouter turns it into a servemux
* But then what is it before you call `Subrouter`?
* In any event, we can now delete the __entire__ `ServeHTTP` method we wrote because I guess gorilla is dealing with all of that now.
* We can now use all the handler methods we wrote (e.g. `getProducts`) and pass them right into the `HandleFunc` method in gorilla
  * we have to convert them to public functions by making them pascal case
  * `getProducts` -> `GetProducts`
* Here's we we are now

```go
ph := handlers.NewProducts(l)

sm := mux.NewRouter()

getRouter := sm.Methods("GET").Subrouter()
getRouter.HandleFunc("/", ph.GetProducts)

s := http.Server{
	Addr:         *bindAddress,
	Handler:      sm,
	IdleTimeout:  120 * time.Second,
	ReadTimeout:  1 * time.Second,
	WriteTimeout: 1 * time.Second,
}
```

* all of this is a mux, which is a router, so we still have to call `http.Serve`
* For our PUT route, we're going to have to pull out an id from the route
  * to define a path variable, use curly braces like `api/route/{pathVar}`
  * if you want to define what that variable looks like, you can use a regexp like `api/route/{id:[0-9]+}`
    * this one means the var is called `id` and it's number with digits 0-9 and there are 1 or more digits
* We no longer need to pass an id into `UpdateProduct` since gorilla attach the path variable into a collection called `mux.Vars` and you pass it the request
  * `vars := mux.Vars(r)`
* `vars` has the id we want

```go
// main.go
ph := handlers.NewProducts(l)

sm := mux.NewRouter()

putRouter := sm.Methods(http.MethodPut).Subrouter()
putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)

// handlers/product.go
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	}
	// ...
}
```

* Okay now POST!!!!!
* Now our code is more specific, and clear, and nice, and refactored!

# Body Parsing Middleware
* 
