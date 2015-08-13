package ldpc

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

const contDrillSPARQL = `
CONSTRUCT {
  ?href a ?type  .
  ?href  rdfs:label ?label .
        ?href geo:geometry ?geometry .
}
FROM <http://data.geolink.org/id/iodp>
WHERE {  
    ?href a ?type ;                         
    rdfs:label ?label;                      
    geo:geometry ?geometry .                                        
} 
`

func ContainerDrillSites(request *restful.Request, response *restful.Response) {
	sprReturn := ContainerDrillSitesCall()
	if strings.Contains(sprReturn, "# Empty NT") {
		response.WriteHeader(http.StatusNoContent)
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

func ContainerDrillSitesCall() string {
	var buff = bytes.NewBufferString("")
	t, err := template.New("rect_template").Parse(contDrillSPARQL)
	if err != nil {
		log.Printf("container for drillsites template creation failed: %s", err)
	}
	err = t.Execute(buff, "")
	if err != nil {
		log.Printf("container for drillsites template execution failed: %s", err)
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
