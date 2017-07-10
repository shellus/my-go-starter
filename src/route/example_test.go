package route_test

import (
	"testing"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/shellus/pkg/logs"
)

func TestExample(t *testing.T) {
	r := mux.NewRouter()
	r.NewRoute().Path("/a/{key}")
	r.NewRoute().Path("/b")

	var match mux.RouteMatch

	req, err := http.NewRequest("GET", "/a/123", nil)
	if err != nil {
		panic(err)
	}

	ok := r.Match(req, &match)
	logs.Debug(ok)
	logs.Debug(match.Vars["key"])
}