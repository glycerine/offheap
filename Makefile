all: gen
	go build
	go install

gen:
	go generate

test: gen
	go test

clean:
	rm -f *~

doc:
	godoc $$GOPATH/src/github.com/remerge/offheap
