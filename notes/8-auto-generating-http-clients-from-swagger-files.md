# Intro
* How can we use the swagger docs we made in episode 7 to generate a client?
* Nic claims this will be a shorter video but it's 46 minutes
* Recall that last time we generated all of our docs and served it at an endpoint
  * `localhost:9090/docs`
* we're gonna use go-swagger again for client get, I think? Kinda spaced.
* I'm very curious how this is going to work out.
* Will there be blanks in the generated code for how to implement each handler?

# Doing the thing
* `swagger generate client` plus options
* Going to make a new folder `client`
* this folder will be a sibling of the `product-api` folder
* While in this folder, we can call

```bash
> swagger generate client -f ../product-api/swagger.yaml -A product-api
```

* Oops when we call that we get an error because the swagger doc is incomplete

```bash
> swagger generate client -f ../product-api/swagger.yaml product-api
2022/11/11 13:20:40 validating spec /Users/mkoshako/code/micro-go/product-api/swagger.yaml
The swagger spec at "/Users/mkoshako/code/micro-go/product-api/swagger.yaml" is invalid against swagger specification 2.0. see errors :
- path param "id" is not present in path "/products"
- path param "{id}" has no parameter definition
```

* Nic's translation: we have a parameter defined in a route but not in the docs
* Specifically here

```go
// swagger:route GET /products/{id} products listSingle
// Return a list of products from the database
// responses:
// 		200: productResponse
// 		404: errorResponse
```

* We actually define it in the `docs.go` file, which I guess isn't good enough
* So we've defined the parameter but haven't linked it to the route (???)
* Where it says `listSingle` above has to match the tag we used in the ID param in `docs.go`
* New

```go
// swagger:route GET /products/{id} products listSingleProduct
// Return a list of products from the database
// responses:
// 		200: productResponse
// 		404: errorResponse

// swagger:parameters listSingleProduct deleteProduct
type productIDParamsWrapper struct {
	// The id of the product for which the operation relates
	// in: path
	// required: true
	ID int `json:"id"`
}
```

* we can use multiple tags in the param def to link it to multiple endpoints
* On to generating the client
* I had to run `go mod init client` in the `client` directory for the gen to succeed
* Also have to run `go mod tidy`
* Lots of stuff got generated!
* Hmm lot's of _unfamiliar_ stuff
* Gonna write a test to see what's going on
* Oop he moved it into the `product-api` folder
  * I will not be doing this
* I mean my test is failing because I declare `c` but never use it... idk why his is
* bloop bloop he's still debugging
* Okay off we go
* My test is passing
* And now it's failing because of connection refused

```go
func TestOurClient(t *testing.T) {
	c := client.Default
	params := products.NewListProductsParams()
	prods, err := c.Products.ListProducts(params)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(prods)
}
```

* Nic is trying to figure out how to create a client with a specific host port
* I just realized that, in my head, I've been thinking we were generating a server, but we're not. 
* Assuming a server is up and running somewhere, where generating a CLIENT - that is a consumer of the API
* Okay so anyway this is what the test looks like now

```go
func TestOurClient(t *testing.T) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)

	params := products.NewListProductsParams()
	prods, err := c.Products.ListProducts(params)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(prods)
}
```

* And this is still giving an error but we can see from the logs on the server that it's getting a request
* Error is `main_test.go:18: &[] (*[]*models.Product) is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface`
* 
