package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/tomachalek/mymquery/api"
)

// API authentication methods (used only if API_TOKEN is set)
const (
	// authCNC logs in via the CNC (korpus.cz) toolbar authentication:
	// the token is exchanged for a session cookie at the login URL
	authCNC = "cnc"
	// authHeader sends the token directly to MQuery in an HTTP header
	// (MQuery's auth.TokenHeaderName), e.g. for a direct access within
	// an internal network
	authHeader = "header"
)

// apiAuth is a configuration of the MQuery API authentication
type apiAuth struct {
	method     string
	token      string
	loginURL   string
	headerName string
}

// validate normalizes the method and checks the configuration
func (a *apiAuth) validate() error {
	a.method = strings.ToLower(strings.TrimSpace(a.method))
	a.headerName = strings.TrimSpace(a.headerName)
	switch a.method {
	case "":
		a.method = authCNC
	case authCNC:
	case authHeader:
		if a.token != "" && a.headerName == "" {
			return fmt.Errorf("auth method %q requires a header name (-token-header or API_TOKEN_HEADER)", authHeader)
		}
	default:
		return fmt.Errorf("unknown auth method %q (expected %s or %s)", a.method, authCNC, authHeader)
	}
	return nil
}

// applyHeaders adds the token header to headers (authHeader method only).
// A header set explicitly (-H) takes precedence.
func (a *apiAuth) applyHeaders(headers map[string]string) {
	if a.method != authHeader || a.token == "" {
		return
	}
	for k := range headers {
		if strings.EqualFold(k, a.headerName) {
			return
		}
	}
	headers[a.headerName] = a.token
}

// login performs a login (authCNC method only)
func (a *apiAuth) login(ctx context.Context, client *api.Client) error {
	if a.method != authCNC || a.token == "" {
		return nil
	}
	return client.Login(ctx, a.loginURL, a.token)
}
