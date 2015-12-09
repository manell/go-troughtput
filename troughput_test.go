package troughput

import (
	"github.com/codegangsta/negroni"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTroughput(t *testing.T) {
	tp := NewTroughput()

	n := negroni.New()
	n.Use(tp)
	n.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	}))

	req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	recorder := httptest.NewRecorder()

	n.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatal("Invalid status code")
	}
}
