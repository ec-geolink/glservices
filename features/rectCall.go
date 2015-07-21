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

type rectVars struct {
	Lat0    string
	Long0   string
	Lat1    string
	Long1   string
	Lat2    string
	Long2   string
	Lat3    string
	Long3   string
}

const spatialRectCall = `
CONSTRUCT {
  ?href a ?type  .
  ?href  rdfs:label ?label .
        ?href geo:geometry ?geometry .
}
WHERE { 
    ?href a ?type ;                         
    rdfs:label ?label;                      
    geo:geometry ?geometry .                                        
    FILTER(bif:st_intersects(?geometry, bif:st_geomfromtext("POLYGON(({{.Long0}} {{.Lat0}}, {{.Long1}} {{.Lat1}}, {{.Long2}} {{.Lat2}}, {{.Long3}} {{.Lat3}}))")))                 
} 
`

func RectFeatures(request *restful.Request, response *restful.Response) {
	sprReturn := RectCall(request.QueryParameter("lat0"), request.QueryParameter("long0"), request.QueryParameter("lat1"), request.QueryParameter("long1"),request.QueryParameter("lat2"), request.QueryParameter("long2"),request.QueryParameter("lat3"), request.QueryParameter("long3"))
	dataparsed, _ := jsonld.ParseDataset([]byte(sprReturn))
	jldOptions := jsonld.NewOptions("http://data.oceandrilling.org")
	jsonldResults := jsonld.FromRDF(dataparsed, jldOptions)
	response.WriteEntity(jsonldResults)
}

func RectCall(lat0 string, long0 string,lat1 string, long1 string,lat2 string, long2 string,lat3 string, long3 string) string {
	params := rectVars{Lat0: lat0, Long0: long0,Lat1: lat1, Long1: long1,Lat2: lat2, Long2: long2,Lat3: lat3, Long3: long3}
	log.Printf("Struct %v \n", params)
	var buff = bytes.NewBufferString("")
	t, err := template.New("rect_template").Parse(spatialRectCall)
	if err != nil {
		log.Printf("rect template creation failed: %s", err)
	}
	err = t.Execute(buff, params)
	if err != nil {
		log.Printf("rect template execution failed: %s", err)
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
