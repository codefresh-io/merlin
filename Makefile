.PHONY: build
build: build-local

build-local: build-classes
	sh ./hack/build.sh

build-classes:
	@echo "Compiling ts into js"
	@tsc --out hack/js/classes.js hack/js/classes.ts --target ES5 --lib es2015 -d
	@echo "Generating go code from JS"
	@go generate hack/generate.go
	@echo "Moving generated code into js package at /pkg/js"
	@mv hack/js/templates.go pkg/js/
