package features

import (
	"bytes"
	//"encoding/json"
	"github.com/emicklei/go-restful"
	jsonld "github.com/linkeddata/gojsonld"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/url"
	"net/http"
	"strings"
	"text/template"
)

const wktPolySPARQL = `
CONSTRUCT {
  ?href a ?type  .
  ?href  rdfs:label ?label .
        ?href geo:geometry ?geometry .
}
WHERE {  
    ?href a ?type ;                         
    rdfs:label ?label;                      
    geo:geometry ?geometry .                                        
    FILTER(bif:st_intersects(?geometry, bif:st_geomfromtext("{{.}}")))                 
} 
`

const errorJSON = `{ "success": false,
  "payload": {
  },
  "error": {
    "code": 123,
    "message": "Empty result set"
  }
}
`

// use web arch error codes?  204 No Content 
// Is the best error from the JSON-LD lib really a strings contains check?
func WKTPoly(request *restful.Request, response *restful.Response) {
	sprReturn := WKTPolyCall(request.QueryParameter("g"))
	if strings.Contains(sprReturn, "# Empty NT") {
		// response.mimetype JSON
		response.WriteHeader(http.StatusNoContent)
		//response.WriteEntity(errorJSON)  // not setting to JSON
	} else {
		dataparsed, error := jsonld.ParseDataset([]byte(sprReturn))
		if error != nil {
			response.WriteEntity("Need JSON error / emptry string here for panic situation")
		}
		jldOptions := jsonld.NewOptions("http://data.geolink.org")
		jsonldResults := jsonld.FromRDF(dataparsed, jldOptions)
		response.WriteEntity(jsonldResults)
	}

}

func WKTPolyCall(g string) string {
	var buff = bytes.NewBufferString("")
	t, err := template.New("rect_template").Parse(wktPolySPARQL)
	if err != nil {
		log.Printf("rect template creation failed: %s", err)
	}
	err = t.Execute(buff, g)
	if err != nil {
		log.Printf("rect template execution failed: %s", err)
	}

	url := "http://data.geolink.org:8890/sparql?default-graph-uri=&query=" + url.QueryEscape(string(buff.Bytes()))
	request := gorequest.New()
	resp, body, errs := request.Get(url).Set("Accept", "text/plain").End()
	if errs != nil {
		log.Printf("Response is an error: %s", errs)
	}
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)

	return body
}
