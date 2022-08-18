# Intro
* Going to be doing JSON validation using a tool called [Go validator](https://github.com/go-playground/validator)
* This will help with securely deserializing our data too
* So it validates the input is what we want AND it makes the app more secure!

# Validating the product
* we already have a struct that represents what the data looks like
* we can leverage this stuct using additional struct tags to help validate the data
* Custom validation types
* go validator has a `validate` method that we pass in our struct
* validator lets us set required fields, specify bounds for things like number inputs (e.g. $0.00 \lt price \leq 100.00$), specify what _kind_ of string an input is (e.g. email string), etc.
    * seems pretty neat!
* Going to add a validate method on our product
* We're going to register the validator inside the `Validate` method but it can also be defined globally

```go
func (p *Product) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
```

* Now to add the validation tags
* These are done similar to the json tags, but the "key" is `validate`
* e.g.

```go
Name        string  `json:"name" validate:"required"`
```

* Ooh gonna write a simple unit test for this. Let's go!
* Write a simple test that we know will fail by validating an empty product

```go
func TestChecksValidation(t *testing.T) {
	p := &Product{}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
```

* The error message we get from `validator` looks like

```shell
/usr/local/opt/go/libexec/bin/go tool test2json -t /private/var/folders/f0/dzkyxv5d63l06f9h9wr4glph0000gn/T/GoLand/___TestChecksValidation_in_product_api_data.test -test.v -test.paniconexit0 -test.run ^\QTestChecksValidation\E$
=== RUN   TestChecksValidation
    products_test.go:11: Key: 'Product.Name' Error:Field validation for 'Name' failed on the 'required' tag
--- FAIL: TestChecksValidation (0.00s)

FAIL

Process finished with the exit code 1
```

* Nice
* Just add the required fields to the product struct to make the test pass
* validator also lets us do custom validation!
* To do custom validation, you have to define a validator with [`validator.RegisterValidation`](https://pkg.go.dev/github.com/go-playground/validator/v10#Validate.RegisterValidation)
* `func (v *Validate) RegisterValidation(tag string, fn Func, callValidationEvenIfNull ...bool) error`
* We'll validate the SKU using a regexp
    * _I should probably read up on go regexps at some point_
* Here's our custom validator

```go
func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}
	
	return true
}
```

* First of all, we're saying our sku needs to look like three strings separated by dashes
* If there isn't only one match (0 or more than 1), validation fails
* To use the custom validator, we first have to pass the func into `RegisterValidation` with a tag name, then we have to use that tag in the struct tag

```go
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"` // <--- HERE!
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)

	return validate.Struct(p)
}
```

# Integrating validation with the API
* recall we've already defined validation middleware for put and post that makes sure the JSON body can be converted into a struct
* We can do our validation there!

```go
func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// ...

		// validate the product
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(
				rw,
				fmt.SprintF("Error validating product: %s", err),
				http.StatusBadRequest
			)
			return
		}
		
		// ...
	})
}
```

* Okay now we test
* First let's make a call that fails validation

```bash
$ curl localhost:9090 -XPOST -d '{"Name":"New Product"}'                                         > Error validating product: Key: 'Product.Price' Error:Field validation for 'Price' failed on the 'gt' tag
> Key: 'Product.SKU' Error:Field validation for 'SKU' failed on the 'required' tag
```

* Very nice error message
* A passing call has no response so I won't waste screen space with it :)