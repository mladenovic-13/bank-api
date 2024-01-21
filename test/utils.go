package test

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/mladenovic-13/bank-api/app"
)

type TestServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *TestServer {
	ts := httptest.NewServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	return &TestServer{ts}
}

func RunTestServer(t *testing.T) *TestServer {
	err := godotenv.Load("../../.env")

	url := os.Getenv("DB_URL")

	if url == "" {
		t.Fatal("failed to load DB_URL env")
	}

	if err != nil {
		t.Fatal(err)
	}

	router, err := app.SetupServer(url)

	if err != nil {
		t.Fatal(err)
	}

	ts := newTestServer(t, router)

	return ts
}

func (ts *TestServer) Get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}

func (ts *TestServer) Post(t *testing.T, urlPath string, body io.Reader) (int, http.Header, []byte) {
	rs, err := ts.Client().Post(ts.URL+urlPath, "application/json", body)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	responseBody, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, responseBody
}
