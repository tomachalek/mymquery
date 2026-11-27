// Package mqmcp provides MQuery API operations as MCP tools (the AI chat
// uses them in-process).
//
// It is a copy of the mcp package of MQuery (github.com/czcorpus/mquery,
// mcp/, as of commit b11e850), as MQuery's module path ("mquery") cannot be
// imported from another module. Changes of the tools should go to MQuery
// as well (it uses them in its MCP server); once MQuery has an importable
// module path, this copy should be replaced by an import.
package mqmcp

import (
	"net/http"

	"github.com/czcorpus/cnc-gokit/logging"
)

type ServerConf struct {
	ListenAddress string `json:"listenAddress"`
}

type Conf struct {
	APIUrl     string              `json:"apiUrl"`
	APIHeaders map[string]string   `json:"apiHeaders"`
	Server     ServerConf          `json:"server"`
	Logging    logging.LoggingConf `json:"logging"`

	// HTTPClient is used for API requests (nil = http.DefaultClient).
	// It allows an embedding application to share e.g. a cookie jar
	// with an authenticated session.
	HTTPClient *http.Client `json:"-"`
}

func (conf *Conf) httpClient() *http.Client {
	if conf.HTTPClient != nil {
		return conf.HTTPClient
	}
	return http.DefaultClient
}
