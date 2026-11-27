package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tomachalek/mymquery/api"
)

func TestAPIAuthValidate(t *testing.T) {
	cases := []struct {
		name    string
		auth    apiAuth
		method  string
		wantErr bool
	}{
		{"default is cnc", apiAuth{}, authCNC, false},
		{"case and spaces", apiAuth{method: " Header ", token: "t", headerName: "X-Tok"}, authHeader, false},
		{"header without name", apiAuth{method: authHeader, token: "t"}, "", true},
		{"header without token", apiAuth{method: authHeader}, authHeader, false},
		{"unknown", apiAuth{method: "basic"}, "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.auth.validate()
			if (err != nil) != c.wantErr {
				t.Fatalf("validate() error = %v, wantErr %v", err, c.wantErr)
			}
			if err == nil && c.auth.method != c.method {
				t.Errorf("method = %q, want %q", c.auth.method, c.method)
			}
		})
	}
}

func TestAPIAuthApplyHeaders(t *testing.T) {
	a := apiAuth{method: authHeader, token: "secret", headerName: "X-Api-Key"}
	h := map[string]string{"Other": "1"}
	a.applyHeaders(h)
	if h["X-Api-Key"] != "secret" || h["Other"] != "1" {
		t.Errorf("unexpected headers %v", h)
	}

	// an explicit -H header wins (case-insensitive)
	h = map[string]string{"x-api-key": "explicit"}
	a.applyHeaders(h)
	if len(h) != 1 || h["x-api-key"] != "explicit" {
		t.Errorf("explicit header overridden: %v", h)
	}

	// cnc method does not touch headers
	h = map[string]string{}
	(&apiAuth{method: authCNC, token: "secret", headerName: "X-Api-Key"}).applyHeaders(h)
	if len(h) != 0 {
		t.Errorf("cnc method added headers: %v", h)
	}
}

func TestAPIAuthHeaderSentToServer(t *testing.T) {
	var got string
	logins := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			logins++
		}
		got = r.Header.Get("X-Api-Key")
		w.Write([]byte(`{"name":"mquery"}`))
	}))
	defer srv.Close()

	a := apiAuth{method: authHeader, token: "secret", headerName: "X-Api-Key", loginURL: srv.URL}
	if err := a.validate(); err != nil {
		t.Fatal(err)
	}
	headers := map[string]string{}
	a.applyHeaders(headers)
	client, err := api.NewClient(srv.URL, headers, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.login(context.Background(), client); err != nil {
		t.Fatal(err)
	}
	if logins != 0 {
		t.Errorf("header method performed a login")
	}
	if _, err := client.ServerInfo(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got != "secret" {
		t.Errorf("token header = %q, want %q", got, "secret")
	}
}
