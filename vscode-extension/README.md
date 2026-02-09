# PolyGo VS Code Extension

Lightweight syntax highlighting and keyword completions for `.p.go` files.

## Locale resolution

The extension reads the locale from `.pgo_lang` at the workspace root. Set it with:

```bash
pgo set jp
```

You can override it via VS Code settings:

- `polygo.locale`: e.g. `bn`, `es`, `jp`, `zh`

## Install (Option B: VSIX)

```bash
cd vscode-extension
npm install
npm run compile
npm i -g @vscode/vsce
vsce package
code --install-extension polygo-vscode-0.0.1.vsix
```

## Dev Mode

```bash
code --extensionDevelopmentPath ./vscode-extension
```

Open a `.p.go` file in the Extension Development Host to verify.
