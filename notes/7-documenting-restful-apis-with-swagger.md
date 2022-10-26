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
