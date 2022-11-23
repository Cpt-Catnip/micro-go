# Intro
* Hey and we're back!
* I'm well, thanks
* Managing CORS requests
* Seems Nic has made a bit of a frontend
  * I really hope I'm not expected to do that too
  * If I am, I'm just going to clone what he did
* The React frontend can't get product data because it doesn't have access to `localhost:9090` origin
* Although annoying, this is for our own protection
* So why can my shell get data?

# How to add CORS control
* Gorilla has a CORS middleware
* Nic is using `yarn` on the frontend...
* Okay whatever I have the frontend clones and running
  * Doing it with Git was actually pretty cool
  * See [this StackOverflow thread](https://stackoverflow.com/a/1355990/12788499)
* Okay back to the API
* We can create a CORS handler with specified origins then wrap our servemux in it so all requests have the header

```go
// CORS
ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

// create a new server
s := http.Server{
    Addr:         *bindAddress,      // configure the bind address
    Handler:      ch(sm),            // set the default handler    <---- HERE
    ErrorLog:     l,                 // set the logger for the server
    ReadTimeout:  5 * time.Second,   // max time to read request from the client
    WriteTimeout: 10 * time.Second,  // max time to write response to the client
    IdleTimeout:  120 * time.Second, // max time for connection using TCP Keep-Alive
}
```

* I guess that's a thing you can do?
* Okay let's spin up the server and see what happens on the frontend
* Yay, we have stuff!!
* We now have the response header `Access-Control-Allow-Origin: http://localhost:3000`
  * which is where our frontend lives
* Why not just allow every origin with `[]string{"*"}`, asks Nic
* You can't do that if you need any kind of authentication!
* It's just _dang_ insecure
* The backend service must correct the correct headers
* The end!
* Okay NOW the next episode will be serving files
* Then a few more episodes and on to gRPC
  * IDK if I'm going to really try to get through all of gRPC
