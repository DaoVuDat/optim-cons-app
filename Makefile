
.PHONY: build
build:
	go build -o bin/optim-cons

build-linux:
	wails build -platform linux -tags webkit2_41

build-windows:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ CGO_CXXFLAGS="-IC:\msys64\mingw64\include" wails build -ldflags '-extldflags "-static"' -skipbindings