build-convstruct:
	go build -o bin/convstruct ./cmd/convstruct

build-gensql:
	go build -o bin/gensql ./cmd/gensql

build-class100:
	go build -o bin/class100 ./cmd/class100

test-convstruct: build-convstruct
	./bin/convstruct \
	-in ./examples/convstruct/pkg-test1 \
	-structName CustomModel \
	-out ./examples/convstruct/pkg-out/out.go  \
	-outPkgName out \
	-debug=true 