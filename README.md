reverse http proxy
==================

Tired of CORS issues? This is a solution for you! This proxy server allows you to serve your content and someones apis from a single host eliminating need of CORS headers. 

### What is it?

`eyjafjallajokull` is a [reverse proxy](http://en.wikipedia.org/wiki/Reverse_proxy). In few words: it's an http server that diverts requests to other other http servers(`endpoints`). 

### How to use it?

#### 1. Add your endpoints to a config file

Config is a json formatted list of routes. By default `eyjafjallajokull` looks for a `routes.json` in current directory to get routes configuration. 

Route properties are:
- `hostname` property will defines a value for `Host` header
- `endpoint` is a url to a server that will respond to requests
- `path` is a regex pattern for request routing by URI


Let's say:
- Your grunt(or other) server lives at `localhost:9000`
- API that you want to use lives at `derpeddit.herokuapp.com`
- APIs that `derpeddit.herokuapp.com` provides start with a prefix `/api`

Here's how a config will look like:

```json
[
	{	"path":"/api",
		"endpoint":"https://derpeddit.herokuapp.com",
		"hostname": "derpeddit.herokuapp.com"
	},
	{	"path":"/",
		"endpoint":"http://localhost:9000", 
		"hostname":"localhost:9000"
	}
]
````

#### 2. Run proxy

```
./eyjafjallajokul -h
  -listen-address="127.0.0.1": address to bind reverse proxy
  -listen-port="8999": port to bind reverse proxy
  -routes="routes.json": path to config file
./eyjafjallajokul -routes="/path/to/routes" -listen-address="127.0.0.1" -listen-port="8999"

```

Et voila! Access your server at localhost:8999


### Download binaries 

- [Windows x64](https://github.com/superduper/eyjafjallajokull/raw/master/build/eyjafjallajokull.amd64.exe)
- [Windows i386](https://github.com/superduper/eyjafjallajokull/raw/master/build/eyjafjallajokull.i386.exe)
- [Linux x64](https://github.com/superduper/eyjafjallajokull/raw/master/build/eyjafjallajokull.linux.amd64)
- [Linux i386](https://github.com/superduper/eyjafjallajokull/raw/master/build/eyjafjallajokull.linux.i386)
- [Mac OS X x64](https://github.com/superduper/eyjafjallajokull/raw/master/build/eyjafjallajokull.linux.amd64)

