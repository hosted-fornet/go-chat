# go-chat

Chat in Golang.
Interops with Rust.

## Setup & Building

1. Get Golang from https://go.dev/dl/
2. Get TinyGo from https://tinygo.org/getting-started/install/linux/
3. ```bash
   cargo install --git https://github.com/bytecodealliance/wit-bindgen wit-bindgen-cli

   git clone https://github.com/hosted-fornet/py-chat.git

   cd go-chat/src/
   go generate

   tinygo build -target=wasi -o go-chat.wasm go-chat.go
   wasm-tools component embed --world process wit/ go-chat.wasm -o go-chat.embed.wasm
   wasm-tools component new go-chat.embed.wasm --adapt wasi_snapshot_preview1.wasm -o go-chat.component.wasm
   cp go-chat.component.wasm ../pkg/go-chat.wasm
   cd ../
   uqdev start-package --url http://localhost:8080
   ```

## Usage

```
/m our@chat:chat:uqbar {"Send": {"target": "foo.uq", "message": "poggers"}}
```
