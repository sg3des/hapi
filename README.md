# HAPI - is a Human API

HAPI is a universal http handler allows used any functions how http handlers.
Current version has only one wrapper over `echo` web framework, in the future
plans to create interfaces to the many other frameworks.

## INSTALL

```sh
go get github.com/sg3des/hapi
```

## USAGE

```go

func main() {
	e := echo.New()
	// some actions...

	h := hapi.New(e)
	h.POST("/some/path", handler)
}

type reqData struct {
	Something string
	Other int
	// ... any other fields
}

type respData struct {
	ID uint
	Name string
	Text string
	// ... etc ...
}

// handler receivs specified strcture filled by JSON or Form data from request
// returns data convert to the JSON and returned how response
func handler(data yourData) respData {
	// ... something actions

	return respData{ID: 123, Name: "response"}
}

```

HAPI support many cases of incoming arguments: parse whole request to struct,
receive files from multipart form data, multiple arguments can be filled from
form data or url values, etc

If function returns:

    - structs returns how JSON
    - string or bytes write as it is to the response body
    - error returns how error
    - int - transform to the status code
