URLShort
=====================

Easy to use URL shortener written in Go and Javascript, for when you just need to host your own 
private URL shortener.

## Build and Run
### 1. Download Packages
~~~
go get github.com/bmizerany/pat
go get github.com/boltdb/bolt
~~~

### 2. Compile The Program
~~~
go build main.go
~~~

### 3. Run It
Startup the server with `./main`. Options which can be specified are:

~~~
-port <PORT>            Port number to listen on
-data <DB_LOCATION>     Where to place the database
~~~

### 4. Shorten some URLs
No URL shortener would be complete without any short URLs. Visit `http://{shortener_server}/shorten` 
to bring up the admin interface and generate short URLs.

### NB: Authentication
This server currently does not support authentication, which means anyone who knows the link to
the shortener page can use it. I suggest using a reverse proxy such as nginx in order to perform 
authentication.

Also place a custom 404 page while you're at it.

