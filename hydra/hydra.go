package hydra

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	jsonld "github.com/linkeddata/gojsonld"
	"log"
)

func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/hydra").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	service.Route(service.GET("/entry").To(Entry))
	service.Route(service.GET("/entry/datasets").To(DataSets))
	return service
}

func Entry(request *restful.Request, response *restful.Response) {
	mapD := map[string]string{"@context": "/hydra/entry/contexts/EntryPoint.jsonld", "@id": "/hydra/entry/", "@type": "EntryPoint", "datasets": "/hydra/entry/datasets/", "datacatalogs": "/hydra/entry/datacatalogs/"}
	mapB, _ := json.Marshal(mapD)
	fmt.Println(string(mapB))

	data, _ := jsonld.ReadJSON(mapB)

	log.Println(string(mapB))
	response.WriteEntity(data)
}

func DataSets(request *restful.Request, response *restful.Response) {
	data := `{
  	 "@context": "/hydra/api-demo/contexts/Collection.jsonld",
  "@type": "Collection",
  "@id": "http://data.geolink.org/hydra/entry/datasets/",
  "members": [
    {
      "@id": "http://data.geolink.org/docs/ds/1",
      "@type": "schema.org:Dataset"
    },
	{
      "@id": "http://data.geolink.org/docs/ds/2",
      "@type": "schema.org:Dataset"
    },
    {
      "@id": "http://data.geolink.org/docs/ds/3",
      "@type": "schema.org:Dataset"
    }
  ]
}`

	dataparsed, _ := jsonld.ReadJSON([]byte(data))
	response.WriteEntity(dataparsed)

}
