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
* Looks like the syntax is awful though- it all happens in comments?
* How do we know where to put the docs?
* Nic is putting it in the handlers directory
* OMG this really is awful.
* I think we're just writting yaml in many single-line comments
    * SURELY there's a better way
* Using `Makefile` to actually do the doc generation
* Yeah I don't think this thing is really working?
    * I can't get it via `brew` or `go get`
    * Might have to skip these next two vids :(
