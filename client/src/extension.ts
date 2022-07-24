/* --------------------------------------------------------------------------------------------
 * Copyright (c) Microsoft Corporation. All rights reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 * ------------------------------------------------------------------------------------------ */

import * as path from 'path';
import { workspace, ExtensionContext } from 'vscode';
import vscode = require('vscode');
import util = require('util');
import cp = require('child_process')

import {
	LanguageClient,
	LanguageClientOptions,
	ServerOptions,
	TransportKind
} from 'vscode-languageclient/node';

let client: LanguageClient;

export async function activate(context: ExtensionContext) {
	// await vscode.window.showErrorMessage("message", "Install");
	const execFile = util.promisify(cp.execFile);
	await execFile("go", ['install', 'github.com/fuji8/gotypeconverter-extension/gotypeconverter-langserver@latest']);

	const cmd = 'gotypeconverter-langserver';
	// const cmd = '/home/fuji/workspace/lsp/gotypeconverter-extension/gotypeconverter-langserver/server';

	// If the extension is launched in debug mode then the debug server options are used
	// Otherwise the run options are used
	const serverOptions: ServerOptions = {
		run: { 
			command: cmd, // wip
		},
		debug: {
			command: cmd,
		}
	};

	// Options to control the language client
	const clientOptions: LanguageClientOptions = {
		// Register the server for plain text documents
		documentSelector: [{ language:  'go'}],
		synchronize: {
			// Notify the server about file changes to '.clientrc files contained in the workspace
			fileEvents: workspace.createFileSystemWatcher('**/.clientrc')
		}
	};

	// Create the language client and start the client.
	client = new LanguageClient(
		'gotypeconverter',
		'gotypeconverter Server',
		serverOptions,
		clientOptions,
		true,
	);

	// Start the client. This will also launch the server
	client.start();
}

export function deactivate(): Thenable<void> | undefined {
	if (!client) {
		return undefined;
	}
	return client.stop();
}
