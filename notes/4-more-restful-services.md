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

# PUTing data
* typically you "put" the entire object
* if you only want to update a few fields, you usually use PATCH
* Okay this introduces another new concetp; usually you specify what you're updating with an ID, so we need to be able to pull the ID used in the request from the URI
  * a bit tricky, actually
* Hmm we're _expecting_ the ID to be in the URI instead of setting up a path like `products/{id}`
* pulling the ID out of a path is __not__ handled by the Go standard library, so you'd need a library/framework like gorilla or gin to do that :(
* nonetheless, we will try
* go regexp (AAAAHHHHHHH)
* `` r := `/([0-9]+)` ``: a `/` then 1 or more digits
* We can use [`MustCompile`](https://pkg.go.dev/regexp@go1.18.4#MustCompile) to make sure that the path in PUT can match this pattern (has an ID)
* Lot to unpack here

```go
if r.Method == http.MethodPut {
	// expect id in the URI
	reg := regexp.MustCompile(`/([0-9]+)`)
	g := reg.FindAllStringSubmatch(r.URL.Path, -1)

	if len(g) != 1 {
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
		return
	}

	if len(g[0]) != 2 {
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
		return
	}

	idString := g[0][1]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
		return
	}

	p.l.Println("got id", id)
}
```

* [regexp.FindAllStringSubmatch](https://pkg.go.dev/regexp@go1.18.4#Regexp.FindAllStringSubmatch)
  * `func (re *Regexp) FindAllStringSubmatch(s string, n int) [][]string`
  * This will a slice of string slices
  * at the first layer, you'll have a list of all strings that match the regexp
    * for our path to be valid then we need only one string to match, otherwise we have something like `/123/42`
  * The second layer will be all _sub_matches, so for us that would be the id itself (__as a string__).
    * IDK what scenario would give us more than one item in this slice, maybe something like `/12ab24` but nonetheless we have to make sure we only move iff there's ~~one item~~ two items in this slice as well
      * __the first item will always be the entire string!!!__ One match means two items.
* The `MustCompile` bit means that if the string we're trying to match _can't_, then we need to panic because, again, the path is invalid.
* we want `g[0][1]` because `g[0][0]` is the entire string that matched, e.g. `/123`
  * we only want `123`
* This is a _lot_ of code just to get an ID from a URI!
* Okay getting and updating the product was kind of a pain in the ass
* First we needed a data method to retrieve the product to update

```go
func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
}
```

* then we needed a data method to update a field on said product and put it back into the product list

```go
func UpdateProducts(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[pos] = p

	return nil
}
```

* This is accepting an id for the product to update and a new product struct to replace it with, then we update the struct at the found position in the list
* then we write a method for the handler to use that unmarshals the JSON and puts it in the product list (uses the new data method)

```go
func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProducts(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
```

* I think we're starting to see how we can really benefit from a framework

# Next time
* We're going to refactor the code using the [Gorilla web toolkit](https://www.gorillatoolkit.org/)
* This is the framework Nic used when he initially started learning Go __9 years ago__
  * I'm sure there's a newer and better one now
  * Nic does still use it though
