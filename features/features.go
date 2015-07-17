package features

import (
	"bytes"
	//"encoding/json"
	"github.com/emicklei/go-restful"
	jsonld "github.com/linkeddata/gojsonld"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/url"
	"text/template"
)

const geolinkSPARQL = `
PREFIX iodp: <http://data.oceandrilling.org/core/iodp/> 
CONSTRUCT {
     ?uri  a iodp:Drillsite  .
	?uri  geo:lat ?lat .
	?uri geo:long ?long .
}  WHERE { 
	?uri  a iodp:Drillsite  .
	?uri  geo:lat ?lat .
	?uri geo:long ?long .
} 
`

const geolinkContext = `
{
  "@context": {
    "glview":	"http://schema.geolink.org/dev/view#",
    "award":	"http://data.geolink.org/id/award/nsf/", 
    "person":	"http://data.geolink.org/id/person/",
    "ical": "http://www.w3.org/2002/12/cal/ical#",
    "xsd": "http://www.w3.org/2001/XMLSchema#",
    "ical:dtstart": {
      "@type": "xsd:dateTime"
    }
  }
}
`

func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/features").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	service.Route(service.GET("/jsonld").To(Features))

	return service
}


func Features(request *restful.Request, response *restful.Response) {

	sprReturn := glSparqlCall("foo")
	dataparsed, _ := jsonld.ParseDataset([]byte(sprReturn))
	jldOptions := jsonld.NewOptions("http://data.geolink.org")
	jsonldResults := jsonld.FromRDF(dataparsed, jldOptions)

	// jldConext := jsonld.Context{context: geolinkContext}
	// jsonldCompacted, _ := jsonld.Compact(dataparsed, jldConext, jldOptions)
	// log.Printf("JSON-LD compacted:\n %s\n\n", jsonldCompacted)

	response.WriteEntity(jsonldResults)
}

func glSparqlCall(query string) string {
	var buff = bytes.NewBufferString("")
	t, err := template.New("lsh_template").Parse(geolinkSPARQL)
	if err != nil {
		log.Printf("sample template creation failed: %s", err)
	}
	err = t.Execute(buff, query)
	if err != nil {
		log.Printf("sample template execution failed: %s", err)
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
