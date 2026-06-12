package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func newTestApplication(t *testing.T) *Application {
	return &Application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

type TestServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *TestServer {
	tlsServer := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	tlsServer.Client().Jar = jar

	tlsServer.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &TestServer{tlsServer}
}

func (ts *TestServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	response, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return response.StatusCode, response.Header, string(body)
}
