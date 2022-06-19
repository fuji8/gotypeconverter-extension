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

func ref[T any](t T) *T {
	return &t
}

func main() {
	// This increases logging verbosity (optional)
	logging.Configure(2, nil)

	handler = protocol.Handler{
		Initialize:  initialize,
		Initialized: initialized,
		Shutdown:    shutdown,
		SetTrace:    setTrace,
		TextDocumentCodeAction: func(context *glsp.Context, params *protocol.CodeActionParams) (interface{}, error) {
			if len(params.Context.Only) >= 1 && params.Context.Only[0] == "refactor" {
				res := &protocol.CodeAction{
					Title: "CodeAction-Test",
					Kind:  ref(protocol.CodeActionKindRefactor),
					Edit: &protocol.WorkspaceEdit{
						Changes: map[protocol.DocumentUri][]protocol.TextEdit{
							params.TextDocument.URI: {
								{
									Range: protocol.Range{
										Start: protocol.Position{
											Line: 1,
										},
										End: protocol.Position{
											Line: 3,
										},
									},
									NewText: "hogeeee",
								},
							},
						},
					},
				}
				return []interface{}{res}, nil
			}
			return nil, nil
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
