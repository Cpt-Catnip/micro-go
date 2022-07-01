recall:
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

# Refactor Handlers // Move away from default server
- gonna refactor code using better practices
- there's too much code in our func main right now, which is bad for testing!
- want contents of `HandleFunc` into its own object
- making a `handlers` package
- `HandleFunc` takes the function you pass into it and it turns it into a handler, which is an interface
```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```
- we want to make a struct that implements this interface
```go
package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Println(("Hello World"))
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(rw, "Hello %s", d)
}
```
- we're using a logger right now, but if we want control over the logging we might want to change that
- Nic mentioned something about not wanting to "create object" inside the handler and needing to use dependency injection for testing purposes
  - we'll get back to that
  - Okay this series is much more confusing to me, which is good and bad
  - [Dependency Injection](https://stackoverflow.com/a/140655)
- This is apparently idiomatic go
```go
type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}
```
- struggling to import the `handlers` package
  - I think another weird "being in the wrong project" issue
  - actually just had to close and reopen vscode
- `HandleFunc` puts a handler on the default mux
- the serve mux is itself a handler (???)
- We can make our own serve mux and pass it into `ListenAndServe`
- Wow our code is so clean now!
```go
// main.go
package main

import (
	"log"
	"net/http"
	"os"
	"working/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	http.ListenAndServe(":9090", sm)
}

// handlers/hello.go
package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println(("Hello World"))
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(rw, "Hello %s", d)
}

// handlers/goodbye.go
package handlers

import (
	"log"
	"net/http"
)

type Goodbye struct {
	l *log.Logger
}

func NewGoodbye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

func (g *Goodbye) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Byeee"))
}
```
- reflect
  - It appears that in go, you make a file for a certain struct. This is an analog for defining a class
  - since there are no classes in proper, the convention is to create a method `NewClass` which returns _a pointer_ to a new struct
    - I imagine the pointer convention just makes it easier and more simple to operate on that struct (dependency injection)
  - It also seems that you almost _never_ pass structs into functions or make methods on structs, always pointers to structs
    - This can/will get confusing since Go lets you call pointer methods on the structs themselves. This means - and I don't like when languages do this - the engineer has to be very cognisant of what's going on under the hood in each particular scenario
  - There also seems to be an overall convention of writing methods that do things particular to the use case, like writing a response, but also designing it in such a way that you can let the engineer customize that behavior, like making the write response method implement the `Writer` interface so you can pass it into any method that accepts that interface
  - Pretty cool
- Still some improvements to be made
- we want to think about defaults and timeouts
- i can't believe this is all in the stdlib
- To do this we need to make our own server
  - `http.ListenAndServe` 
- `s := &http.Server{}`: I'm getting the impression that we almost never want to capture structs in variable and instead just pass around references
  - in practice, there's really no difference between _"pass around a value and make reference when needed"_ and _"pass around a reference and get value when needed"_
- When making our own server instance, we can tune a lot of parameters, like read and write timeouts
- later on we'll get to setting timeouts on a handler-by-handler basis
- something that's particularly useful is the `IdleTimeout`, which is the timeout for a particular TCP connection
  - this is useful when you have many microservices and you want to maintain connections between them
- a lot of this parameter tuning depends on how your service is used
- We'll get into _how_ you can tune these, as in how do you know which tuning to use for your purposes
- Here's what main looks like now
```go
package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"working/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	s.ListenAndServe()
}
```
- You can clearly see the comparrisons to how we were doing it before

# Graceful Shutdown
- We don't want to shut down the server while someone is in the middle of some transaction
- This is important for reliability patters
- _I really like that Nic is talking about the real-world utility of patterns and packages instead of just why it's "good code"_
```go
tc, _ := context.WithDeadline(context.Background(), 30*time.Second)
s.Shutdown(tc)
```
- This says "if the handlers are still working after 30 seconds, forcefully close the server"
- One problem is that `s.ListenAndServer()` is a blocking call, so we need to put it in a goroutine
- Okay now that the listen method is running in a separate routine, the server will move on to `Shutdown` and immediately close the server
- we can use the [os/signal](https://pkg.go.dev/os/signal@go1.18.3) package register certain signals
  - [signal.Notify](https://pkg.go.dev/os/signal@go1.18.3#Notify)
  - `func Notify(c chan<- os.Signal, sig ...os.Signal)`
  - accepts a write only (I think) channel of type `os.Signal` and any number of signals
- Whenever notify receives _that_ signal, it gets passed into the channel
- I suppose it's up to us to decide what to do once a message is received on the channel
- "Pretty much done" he says but we've been staring at his face instead of the code for the past couple of minutes
- Phew okay there we go
- After spinning up the server, we block the main routine by waiting for one of the interrupt signals on the channel we made
- I'm getting a couple of _warnings_ that Nic isn't
  - I'm being told that the channels used in `signal.Notify` should be buffered (because I'm at most sending two messages?)
    - [We must use a buffered channel or risk missing the signal if we're not ready to receive when the signal is sent.](https://pkg.go.dev/os/signal@go1.18.3#example-Notify)
    - I don't understand
    - [By default channels are _unbuffered_, meaning that they will only accept send (`chan <-`) if there is a corresponding receive (`<- chan`) rady to receive the sent value. _Buffered channels_ accept a limited number of values without a corresponding receiver for those values.](https://gobyexample.com/channel-buffering)
    - Okay so why would there not be a receiver ready... I guess since the signal isn't coming from the program, we just can't 100% guarantee that there will be a reciever. Even if we're 99.9% sure, there's still a non-zero chance something goes wrong.
  - I should use the cancel function called from `context.WithTimeout` to about "context leak"
    - what is that
    - [If you fail to cancel the context, the foroutine with that WithCancel or WithTimeout created will be reatined in memory indefinitely (until the program shuts down). causing a memory leak. [...] It's best practice to use a `defer cancel()` immediately after calling `withCancel()` or `WithTimeout()`](https://stackoverflow.com/a/44394873)
    - I guess we're shutting down the program so I don't need to worry, but I will do it anyway for practice
  - New one! "os.Kill cannot be trapped (did you mean syscall.SIGTERM?)"
    - I guess this isn't something that can be caught. Maybe similar to "you can't catch a signal that your computer burst into flames"
- all we need now is unit tests on our handlers
- also going to start looking at implementing more RESTful interfaces
  - verbs, parameters, etc.
  - We'll start with the stdlib and then look at frameworks

```go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"working/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, c := context.WithTimeout(context.Background(), 30*time.Second)
	defer c()
	s.Shutdown(tc)
}
```
