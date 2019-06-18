build-convstruct:
	go build -o bin/convstruct ./cmd/convstruct

build-gensql:
	go build -o bin/gensql ./cmd/gensql

build-class100:
	go build -o bin/class100 ./cmd/class100

test-convstruct: build-convstruct
	./bin/convstruct -in ./convstruct/testdata -structName CustomModel -out ./out.go -outPkgName tools