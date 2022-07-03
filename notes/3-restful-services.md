# Intro
- How to build RESTful services with Go

# REST
- **Re**presentational **S**tate **T**ransfer
- Proposed by Roy Fielding
- One of the most commonly used patters in consumer facing services
  - gRPC is another story and we'll discuss it
- JSON > HTTP
- Specific ways to structure resources
- Not obligated to use JSON, but it's very common to do so

# Make the service RESTful
- oop gotta update my code to how it appears on github :/
- gonna try and build a real world application now!
  - online coffee shop (classic)
- will allow us to demostrate some real world applications
- Ooh there's a lot of new code that happened in the margins
- Not going to use a DB in this episode
- Want to make a GET request to products and return the product list (defined statically for now)
- Can use `encoding/json` to deal with JSON data
- There are two main ways to use this package
- Gonna add a method to the data package that acts as a data access model
  - we want to abstract the data fetching logic away

```go
// data/products.go
func GetProducts() []*Product {
	return productList
}

// handlers/product.go
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	d, err := json.Marshal(lp)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

	rw.Write(d)
}
```

- Logic for getting data is quite trivial right now since we're returning a static list
- Also we're taking advantage of the ResponseWriter `Write` methods; dead easy!
- We want to format the data a little bit before passing to user
  - maybe the DB format isn't user/client friendly
- the json package utilizes [struct tags](https://pkg.go.dev/reflect#StructTag)
  - a struct tag is just an annotation to a field- some packages/code will utilize these
  - json has specific syntax for doing things like ignoring fields
```go
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}
```
- No to tidy up a bit
- [json.NewEndoder](https://pkg.go.dev/encoding/json@go1.18.3#NewEncoder)
  - `func NewEncoder(w io.Writer) *Encoder`
  - returns an encoder that writes to `w`
  - _accepts_ an io writer, so probably has some write methods
  - `rw` is an io writer? (yes)
- yeah so [json.Encode](https://pkg.go.dev/encoding/json@go1.18.3#Encoder.Encode) converts passed data to json and writes it to the writer used in NewEncoder
- What's the benefit here?
- This way we're not buffering anything into memory
- Encoder is also faster than using marshal
  - This sort of thing is imortant when doing multi-thread stuff
- We're going to define a Product slice (`[]*Product`) type so we can define a `toJSON` method on it
- Make new encoder and return the error from the encoding on "self"
```go
func (p *Products) toJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}
```
- also have to update the handler to use new method
- Now reading JSON
- Going to refactor all this resopnding with JSON stuff into its own itnernal (camelCase instead of PascalCase) method since this is all for GET requests only
- recall (mike) that a handler needs to have the `serveHTTP` method to be a valid handler
- this is where ou verb (GET, POST, etc) logic needs to happen
- new `serveHTTP`
```go
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}
```
- now when we make a GET request, we get the data
- on any other type of request, we get a 405: method not allowed response
- In the next episode, we'll learn how to handle an update/PUT request and how to send/receive JSON from the server
