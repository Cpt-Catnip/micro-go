# Recap
* last time we build a microservice with the Go standard library
* helps us appreciate just how much the std lib lets us do
* He's back to a different repo again, this is quite annoying
* Oop no we're back

# Moving along with the product handler
* now we want to introduce methods to update data
  * POST for a new resource
  * PUT for creating or replacing a resource
* Nic says use POST for new and PUT for update (maybe controversial)

# POSTing data
* Gonna follow the same pattern as the GET handler
* Need to define a new condition in the `ServeHTTP` method for handling POSTs
* just by specifying data using the `-d` flag in curl, it knows to do a POST request
* first we need to parse the request data
  * last time we went struct -> JSON, now we want JSON -> struct
* Use encoding/json to decode
  * [http/json NewDecoder](https://pkg.go.dev/encoding/json@go1.18.3#NewDecoder)
  * `func NewDecoder(r io.Reader) *Decoder`
* Then we'll decode into an interface using [Decoder. Decode](https://pkg.go.dev/encoding/json@go1.18.3#Decoder.Decode)
* make a `fromJSON` method on data structure

```go
func (p *Product) fromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
```

* recall that the request body in an io reader
* `Decode` copies the JSON into a passed struct
* Not everything is read at this point in the method handler, which is why a reader is exposed on the request body
* just gonna store data to the fake database, which is an array/slice

```go
// handlers/product.go
func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarchal json", http.StatusBadRequest)
	}

	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}

// data/products.go
func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}
```

* Great! Now how to update data...
