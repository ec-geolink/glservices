package features

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
		Path("/features").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	// service.Route(service.GET("/jsonld").To(Features))
	service.Route(service.GET("/rect").To(RectFeatures))
	service.Route(service.GET("/wktpoly").To(WKTPoly))
	service.Route(service.GET("/drilling").To(DrillingFeatures))
	// Param(service.QueryParameter("lat", "Lattitude of the point of intereest")).
	// Param(service.QueryParameter("long", "Longitude of the point of intereest")).
	// Param(service.QueryParameter("radius", "Radius to sweep for items of intereest")))

	return service
}
