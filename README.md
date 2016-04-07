# JSONP http middleware for [Gin](https://github.com/gin-gonic/gin)

JSONP is a common technique used to communicate with a JSON-serving Web Service with a
Web browser over cross-domains, in place of a XHR request. There is a lot written about
JSONP out there, but the tl;dr on it is a Javascript http client requesting JSONP
will write a `<script>` tag to the head of a page, with the `src` to an API endpoint,
with the addition of a `callback` (or `jsonp`) query parameter that represents a
randomly-named listener function that will parse the request when it comes back from
the server.

This middleware will work for [Gin](https://github.com/gin-gonic/gin). The code
is small, so go read it, but it just buffers the response from the rest of the chain,
and if its a JSON request with a callback, then it will wrap the response in the callback
function before writing it to the actual response writer.

Any feedback is welcome and appreciated!

The origin code is written by [@pkieltyka](https://github.com/pkieltyka)

Changed for [Gin](https://github.com/gin-gonic/gin) by [@jim3mar](https://github.com/jim3mar)

## Example

```go
package main

import (
	"github.com/gin-gonic/gin"
	jsonp "github.com/jim3mar/ginjsonp"
)

func main() {
	r := gin.New()
	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(jsonp.Handler())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
		"message": "pong",
		})
	})
	r.Run(":8088") // listen and server on 0.0.0.0:8080
}
```

*Output:*

```
$ curl -v "http://localhost:8088/ping"
* About to connect() to localhost port 8088 (#0)
*   Trying ::1...
* Connected to localhost (::1) port 8088 (#0)
> GET /ping HTTP/1.1
> User-Agent: curl/7.29.0
> Host: localhost:8088
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Thu, 07 Apr 2016 02:45:40 GMT
< Content-Length: 19
< 
{"message":"pong"}

$ curl -v "http://localhost:8088/ping?callback=X"
* About to connect() to localhost port 8088 (#0)
*   Trying ::1...
* Connected to localhost (::1) port 8088 (#0)
> GET /ping?callback=X HTTP/1.1
> User-Agent: curl/7.29.0
> Host: localhost:8088
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Length: 121
< Content-Type: application/javascript
< Date: Thu, 07 Apr 2016 02:45:58 GMT
< 
* Connection #0 to host localhost left intact
X({"meta":{"content-length":19,"content-type":"application/json; charset=utf-8","status":200},"data":{"message":"pong"}})#
```

## NOTES

Since JSONP must always respond with a 200, as thats what the browser `<script>`
tag expects, a nice pattern that is also used in the GitHub API is to put the HTTP
response headers in a `"meta"` hash, and the HTTP response body in `"data"`. Like so..

```json
JsonpCallbackFn_abc123etc({
  "meta": {
    "Status": 200,
    "Content-Type": "application/json",
    "Content-Length": "19",
    "etc": "etc"
  },
  "data": { "name": "yummy" }
})
```

## LICENSE

BSD 3
