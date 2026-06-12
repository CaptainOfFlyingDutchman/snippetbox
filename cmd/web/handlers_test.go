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
