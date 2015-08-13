package ldpc

import (
	"github.com/emicklei/go-restful"
)

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
		Path("/ldpc").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
    service.Route(service.GET("/").To(ContainerList))
	  service.Route(service.GET("/iodp").To(ContainerDrillSites))
	return service
}
