# PolyGo (`polygo` / `pgo`)

Write Go programs using your native language keywords, compiled by the official Go toolchain.

PolyGo is a multilingual Go language layer, not a Go fork.

---

## ⚡ Quick Start (1–2 minutes)

### 1️⃣ Requirements

* Go **1.25+**
* Linux / macOS / Windows

```bash
go version
```

---

### 2️⃣ Install PolyGo

```bash
go install github.com/newmizanur/poly-go/cmd/pgo@latest
```

Verify:

```bash
pgo version
```

---

### 3️⃣ Try a language (Bangla / Spanish / Japanese / Chinese)

PolyGo supports multiple languages via keyword maps. Choose one and try:

---

## 🅰️ Example: Bangla (`bn.json`)

Create a project:

```bash
mkdir hello-polygo
cd hello-polygo
go mod init hello-polygo
```

Create `main.p.go`:

```go
প্যাকেজ main

আমদানি "fmt"

ফাংশন main() {
	fmt.Println("হ্যালো PolyGo 👋")
}
```

Run:

```bash
pgo run --lang=bn .
```

Output:

```text
হ্যালো PolyGo 👋
```

---

## 🅱️ Example: Spanish (`es.json`)

Create `main.p.go`:

```go
paquete main

importar "fmt"

funcion main() {
	fmt.Println("Hola PolyGo 👋")
}
```

Run:

```bash
pgo run --lang=es .
```

---

## 🅲 Example: Japanese (`jp.json` – Nadesiko-inspired)

```go
パッケージ main

インポート "fmt"

関数 main() {
	fmt.Println("こんにちは PolyGo 👋")
}
```

---

## 🅳 Example: Chinese (`zh.json` – Wenyan-inspired)

```go
包 main

导入 "fmt"

函数 main() {
	fmt.Println("你好 PolyGo 👋")
}
```

For full examples, see the `examples/` folder.

---

## 📂 Project Rules (important)

* Always run PolyGo from the project root
* You can mix:
  * `.go` (normal Go)
  * `.p.go`

---

## 🌐 Locale selection

```bash
pgo set bn
pgo run --lang=jp ./examples/jp.p.go
```

---

## 🔁 Common Commands

```bash
pgo gen        # generate Go files
pgo run .      # generate + run
pgo test ./... # generate + test
pgo build .    # generate + build
pgo clean      # remove .pgo_gen
```

---

## 🌍 Language Support Model

PolyGo is language-agnostic. Each language is defined by a JSON map:

```
bn.json  → Bangla
es.json  → Spanish
jp.json  → Japanese
zh.json  → Chinese
```

You can add your own language by creating a new map and example. Please open a pull request with:

* `lang/<locale>.json`
* `examples/<locale>.p.go`

Or fork and maintain your own language pack.

---

## 🔐 Using keywords as identifiers (`@` escape)

If a keyword conflicts with a variable name:

```go
@型 := "this is a variable"
fmt.Println(型)
```

`@` means “do not treat this as a keyword”.

---

## 🧩 VS Code extension

Marketplace:
```
https://marketplace.visualstudio.com/items?itemName=mizanur.polygo-vscode
```

Dev mode:

```bash
code --extensionDevelopmentPath ./vscode-extension
```

VSIX install:

```bash
cd vscode-extension
npm install
npm run compile
npm i -g @vscode/vsce
vsce package
code --install-extension polygo-vscode-0.0.1.vsix
```

---

## 🧾 License

MIT (see LICENSE).
