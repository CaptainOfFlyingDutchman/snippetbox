package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestUserSignup(t *testing.T) {
	type TestCase struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}

	app := newTestApplication(t)

	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	_, _, body := testServer.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	// t.Logf("CSRF token is: %q", csrfToken)
	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action=\"/user/signup\" method=\"POST\" novalidate>"
	)

	tests := []TestCase{
		{name: "Valid submission", userName: validName, userEmail: validEmail, userPassword: validPassword, csrfToken: csrfToken, wantCode: http.StatusSeeOther},
		{name: "Invalid CSRF Token", userName: validName, userEmail: validEmail, userPassword: validPassword, csrfToken: "wrongToken", wantCode: http.StatusBadRequest},
		{name: "Empty name", userName: "", userEmail: validEmail, userPassword: validPassword, csrfToken: csrfToken, wantCode: http.StatusUnprocessableEntity, wantFormTag: formTag},
		{name: "Empty email", userName: validName, userEmail: "", userPassword: validPassword, csrfToken: csrfToken, wantCode: http.StatusUnprocessableEntity, wantFormTag: formTag},
		{name: "Empty password", userName: validName, userEmail: validEmail, userPassword: "", csrfToken: csrfToken, wantCode: http.StatusUnprocessableEntity, wantFormTag: formTag},
		{name: "Invalid email", userName: validName, userEmail: "bob@example.", userPassword: validPassword, csrfToken: csrfToken, wantCode: http.StatusUnprocessableEntity, wantFormTag: formTag},
		{name: "Short password", userName: validName, userEmail: validEmail, userPassword: "pa$$", csrfToken: csrfToken, wantCode: http.StatusUnprocessableEntity, wantFormTag: formTag},
		{name: "Duplicate email", userName: validName, userEmail: "dupe@example.com", userPassword: validPassword, csrfToken: csrfToken, wantCode: http.StatusUnprocessableEntity, wantFormTag: formTag},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", test.userName)
			form.Add("email", test.userEmail)
			form.Add("password", test.userPassword)
			form.Add("csrf_token", test.csrfToken)

			code, _, body := testServer.postForm(t, "/user/signup", form)

			assert.Equal(t, code, test.wantCode)

			if test.wantFormTag != "" {
				assert.StringContains(t, body, test.wantFormTag)
			}
		})
	}
}
