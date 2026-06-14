package main

import (
	"bytes"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"snippetbox.manvendrask.com/internal/models/mocks"
)

var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+)">`)

func extractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(matches[1])
}

func newTestApplication(t *testing.T) *Application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.IdleTimeout = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &Application{
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		snippets:       &mocks.SnippetModel{},
		users:          &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
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

func (ts *TestServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	// We can't use http.Client.PostForm() because it doesn't support setting the Referer header
	request, err := http.NewRequest(http.MethodPost, ts.URL+urlPath, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Referer", ts.URL+urlPath)

	response, err := ts.Client().Do(request)
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
