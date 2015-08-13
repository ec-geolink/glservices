package main

import (
	"geolink.org/glservices/people"
	"geolink.org/glservices/ldpc"
    "geolink.org/glservices/hydra"
	"geolink.org/glservices/features"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
)

func main() {
	// Simple sevice for some static pages about the glservice
	serverMuxA := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	serverMuxA.Handle("/", fs)

	go func() {
		http.ListenAndServe("localhost:8081", serverMuxA)
	}()

	// The service API code
	wsContainer := restful.NewContainer()
	// u := UserResource{}
	// u.RegisterTo(wsContainer)

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type"},
		CookiesAllowed: false,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)
	// Add container filter to respond to OPTIONS
	wsContainer.Filter(wsContainer.OPTIONSFilter)

	wsContainer.Add(people.New())    // simple people search, uni-file pattern  (multi-file used in feature)
	wsContainer.Add(features.New())  // simple example for testing WKT polygons
	wsContainer.Add(ldpc.New())      // testing Linked Data Patterns
	wsContainer.Add(hydra.New())     // testing hydra hypermedia approach

	log.Printf("Listening on localhost:6789 (service) and localhost:8081 (static files)")

	server := &http.Server{Addr: ":6789", Handler: wsContainer}
	server.ListenAndServe()

}
