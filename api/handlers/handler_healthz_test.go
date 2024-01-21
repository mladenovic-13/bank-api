package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mladenovic-13/bank-api/test"
)

func TestHandleHealthz(t *testing.T) {
	ts := test.RunTestServer(t)
	defer ts.Close()

	code, _, body := ts.Get(t, "/healthz")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d\n", http.StatusOK, code)
	}

	var stringValue string
	err := json.Unmarshal(body, &stringValue)

	if err != nil {
		t.Fatal(err)
	}

	if stringValue != "ok" {
		t.Errorf("want body to equal %q", "ok")
	}
}
