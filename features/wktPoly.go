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

func WKTPoly(request *restful.Request, response *restful.Response) {
	log.Printf("In WKTPoly")
	sprReturn := WKTPolyCall(request.QueryParameter("g"))
	log.Printf("%s \n", sprReturn)
	dataparsed, _ := jsonld.ParseDataset([]byte(sprReturn))
	jldOptions := jsonld.NewOptions("http://data.geolink.org")
	jsonldResults := jsonld.FromRDF(dataparsed, jldOptions)
	response.WriteEntity(jsonldResults)
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
