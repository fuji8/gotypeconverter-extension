package main

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"github.com/tliron/kutil/logging"

	// Must include a backend implementation. See kutil's logging/ for other options.
	_ "github.com/tliron/kutil/logging/simple"
)

const lsName = "my language"

var version string = "0.0.1"
var handler protocol.Handler
var log logging.Logger

func main() {
	// This increases logging verbosity (optional)
	logging.Configure(1, nil)

	handler = protocol.Handler{
		Initialize:  initialize,
		Initialized: initialized,
		Shutdown:    shutdown,
		SetTrace:    setTrace,
		TextDocumentDidOpen: func(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
			log.Info("didopen")
			s := protocol.DiagnosticSeverityWarning
			hoge := "sample"
			context.Notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
				URI: params.TextDocument.URI,
				Diagnostics: []protocol.Diagnostic{
					{
						Range:    protocol.Range{Start: protocol.Position{Line: 1}, End: protocol.Position{Line: 2}},
						Message:  "Hello",
						Severity: &s,
						Source:   &hoge,
					},
				},
			})
			return nil
		},
	}

	server := server.NewServer(&handler, lsName, false)
	log = server.Log

	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (interface{}, error) {
	capabilities := handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
