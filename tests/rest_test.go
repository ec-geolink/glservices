package tests

import (
	"net/http"
	"testing"

	. "github.com/emicklei/forest"
)

// ref: http://ernestmicklei.com/2015/07/04/testing-your-rest-api-in-go-with-forest/

var github = NewClient("https://api.github.com", new(http.Client))

func TestForestProjectExists(t *testing.T) {
	cfg := NewConfig("/repos/emicklei/{repo}", "forest").Header("Accept", "application/json")
	r := github.GET(t, cfg)
	ExpectStatus(t, r, 200)
}