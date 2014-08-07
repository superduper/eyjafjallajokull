package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mailgun/vulcan"
	"github.com/mailgun/vulcan/endpoint"
	"github.com/mailgun/vulcan/headers"
	"github.com/mailgun/vulcan/loadbalance/roundrobin"
	"github.com/mailgun/vulcan/location/httploc"
	"github.com/mailgun/vulcan/netutils"
	"github.com/mailgun/vulcan/request"

	"github.com/mailgun/vulcan/route/pathroute"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type tLocation struct {
	address, port string
}

func (self *tLocation) bindTo() string {
	return fmt.Sprintf("%s:%s", self.address, self.port)
}

type tConfig struct {
	location   tLocation
	routeTable []tRoute
}

type tRoute struct {
	Path     string `json:"path"`
	Hostname string `json:"hostname"`
	Endpoint string `json:"endpoint"`
}

type tHostRewriter struct {
	Hostname string
}

func (self *tHostRewriter) ProcessRequest(r request.Request) (*http.Response, error) {
	req := r.GetHttpRequest()
	req.Host = self.Hostname
	req.Header.Set("Via", "eyjafjallajokull/1.0")
	// Remove hop-by-hop headers to the backend.  Especially important is "Connection" because we want a persistent
	// connection, regardless of what the client sent to us.
	netutils.RemoveHeaders(headers.HopHeaders, req.Header)
	return nil, nil
}

func (self *tHostRewriter) ProcessResponse(r request.Request, a request.Attempt) {

}

func readConfig() (*tConfig, error) {
	// read config
	var configPath, listenAddress, listenPort string
	flag.StringVar(&configPath, "routes", "routes.json", "path to config file")
	flag.StringVar(&listenAddress, "listen-address", "127.0.0.1", "address to bind reverse proxy")
	flag.StringVar(&listenPort, "listen-port", "8999", "port to bind reverse proxy")
	flag.Parse()

	if file, err := ioutil.ReadFile(configPath); err != nil {
		return nil, err
	} else {
		// parse the config
		var cfg []tRoute
		if err := json.Unmarshal(file, &cfg); err != nil {
			return nil, err
		} else {
			return &tConfig{tLocation{address: listenAddress, port: listenPort}, cfg}, nil
		}
	}

}

func main() {
	if cfg, err := readConfig(); err != nil {
		log.Fatalf("Failed to read routes config file: %s", err.Error())
	} else {
		// define locations and endpoints
		router := pathroute.NewPathRouter()
		for _, v := range cfg.routeTable {
			rr, err := roundrobin.NewRoundRobin()
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
			rr.AddEndpoint(endpoint.MustParseUrl(v.Endpoint))
			location, err := httploc.NewLocation(v.Path, rr)

			chain := location.GetMiddlewareChain()
			chain.Remove(httploc.RewriterId) // default rewrite header middleware
			chain.Add(httploc.RewriterId, -2, &tHostRewriter{v.Hostname})
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
			router.AddLocation(v.Path, location)
			fmt.Printf("- Added route: %+v\n", v)
		}

		// create a proxy
		proxy, err := vulcan.NewProxy(router)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		// Proxy acts as http handler
		server := &http.Server{
			Addr:           cfg.location.bindTo(),
			Handler:        proxy,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		fmt.Println("! Serving at:", cfg.location.bindTo())
		server.ListenAndServe()
	}

}
