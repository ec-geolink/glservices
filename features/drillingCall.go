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

type spatialVars struct {
	Lat    string
	Long   string
	Radius string
}

const spatialOceanDrilling = `
PREFIX  geo:  <http://www.w3.org/2003/01/geo/wgs84_pos#> 
PREFIX iodp: <http://data.oceandrilling.org/core/iodp/>  
CONSTRUCT {
  ?uri a iodp:Drillsite  .
   ?uri  geo:lat ?lat .
        ?uri geo:long ?long .
}
WHERE  { 
        ?uri  geo:geometry  ?geo    . 
        ?uri  geo:lat ?lat .
        ?uri geo:long ?long .
          FILTER ( <bif:st_intersects> ( ?geo, <bif:st_point> ({{.Long}}, {{.Lat}}), {{.Radius}}))
} 
`

func DrillingFeatures(request *restful.Request, response *restful.Response) {
	sprReturn := DrillingCall(request.QueryParameter("lat"), request.QueryParameter("long"), request.QueryParameter("radius"))
	dataparsed, _ := jsonld.ParseDataset([]byte(sprReturn))
	jldOptions := jsonld.NewOptions("http://data.oceandrilling.org")
	jsonldResults := jsonld.FromRDF(dataparsed, jldOptions)
	// todo  resolve how (if I need to) insert context here
	// jldConext := jsonld.Context{context: geolinkContext}
	// jsonldCompacted, _ := jsonld.Compact(dataparsed, jldConext, jldOptions)
	// log.Printf("JSON-LD compacted:\n %s\n\n", jsonldCompacted)
	response.WriteEntity(jsonldResults)
}

func DrillingCall(lat string, long string, radius string) string {
	params := spatialVars{Lat: lat, Long: long, Radius: radius}
	var buff = bytes.NewBufferString("")
	t, err := template.New("lsh_template").Parse(spatialOceanDrilling)
	if err != nil {
		log.Printf("sample template creation failed: %s", err)
	}
	err = t.Execute(buff, params)
	if err != nil {
		log.Printf("sample template execution failed: %s", err)
	}

	url := "http://data.oceandrilling.org:8890/sparql?default-graph-uri=&query=" + url.QueryEscape(string(buff.Bytes()))
	request := gorequest.New()
	resp, body, errs := request.Get(url).Set("Accept", "text/plain").End()
	if errs != nil {
		log.Printf("Response is an error: %s", errs)
	}
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)

	return body
}
