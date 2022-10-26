# Intro
* Docs!!!!
* [Swagger!!!!!](https://swagger.io/)
* Using [GoSwagger](https://github.com/go-swagger/go-swagger)
    * note: it looks like this package (module?) is out of date and has failing build. I will of course use this for learning purposes in this video but I should try and find something else
* Nic says it's good to write docs
    * sure, bud
* We want to look at how we can start using swagger with our _existing_ api
* Specifically, going to automatically generate go code

# Go Swagger
* Will allow us to generate documentation from our API
* There are top level docs we still need to write, like how wormhole is working?
* Looks like the syntax is awful though - it all happens in comments?
* How do we know where to put the docs?
* Nic is putting it in the handlers directory
* OMG this really is awful.
* I think we're just writing yaml in many single-line comments
    * SURELY there's a better way
* Using `Makefile` to actually do the doc generation
* Yeah, I don't think this thing is really working?
    * I can't get it via `brew` or `go get`
    * Might have to skip these next two vids :(
* Arthur helped me install it and it works now! Yippee!!
* `go get` is deprecated as of go 1.18. You have to use `go install` now.
* Now on to documenting the APIs!!
* Hold on just a second, I seem to be missing a number of files...
* Yes, Nic did quite a bit of refactoring in between videos here
  * this is incredibly annoying since I have to copy over all the content from episode 7 and not just what Nic has at the beginning
  * I'm gonna try to do it as the files come up
* BRB gonna do a lot of code editing...
* Getting a little ahead of myself and there's stuff I'm adding in that is DEF gonna be explained in the video
* Done with copy pasta and still there are errors >:(
  * Resolved all except the middleware problem, which will likely be addressed in the video
  * Will comment out `middleware`
* **---> Back to the video now!!! <---**
* `make swagger` is returning a very unhelpful error message
  * `unsupported type "invalid type""`
* Like... where is it?
* **--> Back again lol <--**
* This has been a trip; I'm determined to at least watch the rest of this video today.
* installing go swagger in makefile...
* Adding API elements for each part of the API now
* API docs (just go comments) get included in swagger docs
* Nic is able to generate docs here but I'm still getting that error
* Doc generation is working after removing `GO111MODULE=off` from the `swagger` make command (or whatever)
* go-swagger will automatically create references!
* Set response codes/schemas in comments too
* Something tells me this is just going to be a lot of familiar stuff, especially since I already wrote all this :/
* It's useful to define types in Go just to reference them in the docs, for example

```go
// in handlers/get.go
// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//		200: productsResponse

// in handlers/docs.go
// Data structure representing a single product
// swagger:response productResponse
type productResponseWrapper struct {
// Newly created product
// in: body
Body data.Product
}
```

* Note that the type can be called whatever you want so long as the `swagger:response` comment matches what is used in the return type
* I'm not sure why these wrappers are made. Presumably these types are defined elsewhere for actual function signatures.
* For some reason `in: body` is put inside the struct def
* I'm not seeing that annotation anywhere...
* go-swagger will pull in all type def from Go into the docs!
  * That's actually sick

# Serving Docs on API
* we can add a docs handler using a special middleware
* [readoc](https://github.com/Redocly/redoc) has its own middleware for serving swagger docs witha really nice UI
* actual package is in go-openapi runtime/middleware package [here](https://github.com/go-openapi/runtime/tree/master/middleware)
* `go get "github.com/go-openapi/runtime/middleware"`
  * this is giving me a warning :(
* So far we have this, which won't work yet

```go
// handler for documentation
opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
sh := middleware.Redoc(opts, nil)

getR.Handle("/docs", sh)
```

* Readoc doesn't know where the file is!
* We need to serve the docs on a built-in file server
  * not gonna go too in depth to this here
* `getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))`
* When someone asks for `GET /swawgger.yaml`, use a handler that connects to the server's file system
  * [file server docs](https://pkg.go.dev/net/http@go1.19.2#FileServer)
* The server will look for the given path in the directory specified in the file server
  * e.g. `./swagger.yaml`
* in further vids, we'll look into gzipping, which will speed up network communication

# Back to documenting the API
* you can use curly braces to specify parameter refs, like `DELETE /products/{id}`
* Then we can define `id` elsewhere
* Starting to seem like _this_ video isn't going to catch up with the code I copied as prep
* Not all of these annotations are really working...
* Wait I take it back; check this out

```go
// swagger:parameters deleteProduct
type productIDParameterWrapper struct {
	// The id of the product to delete from the database
	// in: path
	// required: true
	ID int `json:"id"`
}
```

* Okay so the top-most comment is saying "this defined a parameter for the `deleteProduct` operation"
* But how does it know I'm defining the `id` parameter? is it because the name of the struct field is `ID`?
* We can add more rich description using the "model" tag, e.g.

```go
// Product defines the structure of an API product
// swagger:model
type Product struct {
	// the id for the product 
	// 
	//required: false 
	//min: 1 
	ID int `json:"id"`
	
	// ...
}
```

* This all now shows up in the response type in the docs
* Have we done anything to say "this method returns this type" or is that all from the function signature?
* Here we say `productResponse` had in its body `data.Product`
* Then, presumably, the above docs define what will be displayed in the docs site

```go
// Data structure representing a single product
// swagger:response productResponse
type productResponseWrapper struct {
	// Newly created product
	// in: body
	Body data.Product
}
```

Okay yeah I'm a little out of sync. That's okay though we will persevere.
