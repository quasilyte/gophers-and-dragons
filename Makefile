www:
	GOOS=js GOARCH=wasm go build -o www/go.wasm ./wasm/

.PHONY: www
