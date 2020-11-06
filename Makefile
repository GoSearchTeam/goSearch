build:
	go build .

install:
	go install .
test:
	go test -v .
log:
	find ./documents/ -type f | xargs tail -n +1 && echo
