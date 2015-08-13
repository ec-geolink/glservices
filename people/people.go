package people

import (
	"bytes"
	"encoding/json"
	utilities "geolink.org/glservices/utilities"
	"github.com/emicklei/go-restful"
	jsonld "github.com/linkeddata/gojsonld"
	"github.com/parnurzeal/gorequest"
	"log"

	"net/url"
	"text/template"
)

const geolinkSPARQL = `
PREFIX glview:	<http://schema.geolink.org/dev/view#> 
PREFIX award:	<http://data.geolink.org/id/award/nsf/> 
PREFIX person:	<http://data.geolink.org/id/person/> 
CONSTRUCT {
   ?uri glview:matches <http://data.geolink.org/id/person/bf402132-8c19-409e-8dc1-864700efb838> .
   ?uri rdfs:label ?label .
} 	
WHERE { 
	?uri glview:matches  <http://data.geolink.org/id/person/bf402132-8c19-409e-8dc1-864700efb838> .  
        <http://data.geolink.org/id/person/bf402132-8c19-409e-8dc1-864700efb838> rdfs:label ?label .
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
		Path("/people").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	service.Route(service.GET("/{leg}").To(People))
	service.Route(service.GET("/test/{leg}").To(People))
	service.Route(service.GET("/jsonld").To(PeopleLD))
	return service
}

func People(request *restful.Request, response *restful.Response) {
	data := utilities.People("113")
	response.WriteEntity(data)
}

func PeopleLD(request *restful.Request, response *restful.Response) {

	sprReturn := glSparqlCall2("foo")
	log.Printf("Data from sparql call:\n %v \n\n", sprReturn)

	// Now convert this to JSON-LD
	dataparsed, _ := jsonld.ParseDataset([]byte(sprReturn))
	//log.Printf(dataparsed)

	jldOptions := jsonld.NewOptions("http://data.geolink.org")

	jsonldResults := jsonld.FromRDF(dataparsed, jldOptions)
	log.Printf("JSON-LD:\n %s\n\n", jsonldResults)

	// jldConext := jsonld.Context{context: geolinkContext}
	// jsonldCompacted, _ := jsonld.Compact(dataparsed, jldConext, jldOptions)
	// log.Printf("JSON-LD compacted:\n %s\n\n", jsonldCompacted)

	jsonString, _ := json.Marshal(jsonldResults)
	log.Println("JSON string for use in JSON-LD linter")
	log.Println(string(jsonString))

	log.Printf("\nSerialized back to RDF:\n %s \n\n", dataparsed.Serialize())

	data := "place holder response string"
	response.WriteEntity(data)
}

func glSparqlCall2(query string) string {
	// build SPARQL string, we will want string(buff.Bytes())
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
