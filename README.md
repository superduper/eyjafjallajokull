reverse http proxy
==================

Tired of CORS issues? This is a solution for you! This proxy server allows you to serve your content and someones apis from a single host eliminating need of CORS headers. 

### How to use it?

1. Add your endpoints to a config file

```json
[
	{	"path":"/api",
		"endpoint":"https://derpeddit.herokuapp.com",
		"hostname": "derpeddit.herokuapp.com"
	},
	{   "path":"/*",
		"endpoint":"http://localhost:8000", 
		"hostname":"localhost:8000"
	}
]
 }]
````

2. Run proxy

```
./eyjafjallajokul -h
  -listen-address="127.0.0.1": address to bind reverse proxy
  -listen-port="8999": port to bind reverse proxy
  -routes="routes.json": path to config file
./eyjafjallajokul -routes="/path/to/routes" -listen-address="127.0.0.1" -listen-port="8999"

```

Et voila! Access your server at localhost:8999
