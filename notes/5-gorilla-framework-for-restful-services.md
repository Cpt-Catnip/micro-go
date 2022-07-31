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
* Like express, gorilla has a `Use` function
* middleware allows us to chain multiple handlers together for doing things like authentication and CORS

```go
r := mux.NewRouter()
r.Use(middlewareFunc)
```

* `MiddlewareFunc` is a function that accepts an `http.Handler` and returnds an `http.Handler`
  * `type MiddlewareFunc func(http.Handler) http.Handler`
  * [Docs](https://github.com/gorilla/mux#middleware)
* Recall: we want to validate the the JSON sent in a PUT/PATCH request represents a valid product
* We can decode the request body into a product structure using the `FromJSON` method we wrote
* If there's any error, we respond with an error and `return`, which will break out of the middleware chain
* Otherwise, we'll, as Nic puts it, put that product somewhere so we can use it later
* That somwhere will be the good ol' `context`
  * request has a `r.Context()` method
* contexts need keys
* You can use strings as keys but the convention is to use types
* ~~I'm gettting some errors that Nic isn't, but I also think he has bugs in his code that he hasn't caught yet.~~
* We then create a copy of the req that we're going to "call the upstream with" whatever that means.
  * We make a copy using the context we just created
* then we call `next,ServeHTTP` to move on to the next bit of middleware?
* Here's the middleware we've just defined

```go
type KeyProduct struct{}

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
```

* Now to implement the middleware!

```go
func main() {
  // ...

  putRouter := sm.Methods(http.MethodPut).Subrouter()
  putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
  putRouter.Use(ph.MiddlewareValidateProduct)

  postRouter := sm.Methods(http.MethodPost).Subrouter()
  postRouter.HandleFunc("/", ph.AddProduct)
  postRouter.Use(ph.MiddlewareValidateProduct)

  // ...
}
```

* gorilla knows to run middleware before the request handler, so we don't have to call `Use` before `HandleFunc`
* NOW we can use the product we've put on the context in the route handlers!
  * Yippee!!
* When we pull the product from the context, we get an interface, which we can cast to a product

```go
prod := r.Context().Value(KeyProduct{}).(data.Product)
```

# Coming up next
* Validation on our structs!
* I kinda thought that's what we'd do this time but I guess we just introduced middleware here
* You guessed it, gorilla has a package for that
