package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.manvendrask.com/internal/assert"
)

func TestPingUnit(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ping(responseRecorder, request)

	response := responseRecorder.Result()
	assert.Equal(t, response.StatusCode, http.StatusOK)

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	statusCode, _, body := testServer.get(t, "/ping")

	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
	type TestCase struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}

	app := newTestApplication(t)

	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	tests := []TestCase{
		{name: "Valid ID", urlPath: "/snippet/view/1", wantCode: http.StatusOK, wantBody: "An old silent pond..."},
		{name: "Non-existend ID", urlPath: "/snippet/view/2", wantCode: http.StatusNotFound},
		{name: "Negative ID", urlPath: "/snippet/view/-1", wantCode: http.StatusNotFound},
		{name: "Decimal ID", urlPath: "/snippet/view/1.23", wantCode: http.StatusNotFound},
		{name: "String ID", urlPath: "/snippet/view/foo", wantCode: http.StatusNotFound},
		{name: "Empty ID", urlPath: "/snippet/view/", wantCode: http.StatusNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, _, body := testServer.get(t, test.urlPath)

			assert.Equal(t, code, test.wantCode)

			if test.wantBody != "" {
				assert.StringContains(t, body, test.wantBody)
			}
		})
	}
}
